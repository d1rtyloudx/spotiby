package auth

import (
	"context"
	"errors"
	"github.com/d1rtyloudx/spotiby/user-service/internal/config"
	"github.com/d1rtyloudx/spotiby/user-service/internal/domain/model"
	"github.com/d1rtyloudx/spotiby/user-service/internal/dto"
	"github.com/d1rtyloudx/spotiby/user-service/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrUserAlreadyRegistered = errors.New("user is already registered")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrTokenBlacklisted      = errors.New("token already exists in blacklist")
)

type tokenBlacklist interface {
	IsExists(ctx context.Context, jti string) (bool, error)
	Add(ctx context.Context, jti string, ttl time.Duration) error
}

type profileProvider interface {
	GetByCredentialID(ctx context.Context, id string) (model.Profile, error)
}

type credentialStorage interface {
	Create(ctx context.Context, cred model.Credential) (string, error)
	Update(ctx context.Context, cred model.Credential) error
	GetByUsername(ctx context.Context, userName string) (model.Credential, error)
	GetByEmail(ctx context.Context, email string) (model.Credential, error)
	GetByID(ctx context.Context, id string) (model.Credential, error)
}

type Service struct {
	credentialStorage credentialStorage
	profileProvider   profileProvider
	tokenBlacklist    tokenBlacklist
	log               *zap.Logger
	cfg               *config.TokenConfig
}

func New(
	credentialsStorage credentialStorage,
	profileProvider profileProvider,
	tokenBlacklist tokenBlacklist,
	log *zap.Logger,
	cfg *config.TokenConfig,
) *Service {
	return &Service{
		credentialStorage: credentialsStorage,
		profileProvider:   profileProvider,
		tokenBlacklist:    tokenBlacklist,
		log:               log,
		cfg:               cfg,
	}
}

func (s *Service) Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	childLog := s.log.With(
		zap.String("op", "auth.Service.Register"),
		zap.String("email", req.Email),
		zap.String("username", req.Username),
	)

	hashPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		childLog.Error("failed to hash password", zap.Error(err))
		return dto.RegisterResponse{}, err
	}

	id, err := s.credentialStorage.Create(ctx, model.Credential{
		Username: req.Username,
		Email:    req.Email,
		HashPass: hashPass,
	})
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return dto.RegisterResponse{}, ErrUserAlreadyRegistered
		}
		childLog.Error("failed to register credential", zap.Error(err))
		return dto.RegisterResponse{}, err
	}

	childLog.Info("successfully registered")

	return dto.RegisterResponse{
		ID: id,
	}, nil
}

func (s *Service) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	childLog := s.log.With(
		zap.String("op", "auth.Service.Login"),
		zap.String("username", req.Username),
	)

	cred, err := s.credentialStorage.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return dto.LoginResponse{}, ErrInvalidCredentials
		}
		childLog.Error("failed to get credential by username", zap.Error(err))
		return dto.LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword(cred.HashPass, []byte(req.Password))
	if err != nil {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	profile, err := s.profileProvider.GetByCredentialID(ctx, cred.ID)
	if err != nil {
		childLog.Error("failed to get user profile", zap.Error(err))
		return dto.LoginResponse{}, err
	}

	tokenPair, err := s.newTokenPair(cred, profile.ID)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	childLog.Info("successfully login")

	return dto.LoginResponse{
		Profile: dto.Profile{
			ID:           profile.ID,
			FirstName:    profile.FirstName,
			LastName:     profile.LastName,
			Description:  profile.Description,
			AvatarURL:    profile.AvatarURL,
			CredentialID: cred.ID,
		},
		Tokens: tokenPair,
	}, nil
}

func (s *Service) UpdatePassword(ctx context.Context, id string, password string) error {
	childLog := s.log.With(zap.String("op", "auth.Service.UpdatePassword"))

	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		childLog.Error("failed to hash password", zap.Error(err))
		return err
	}

	err = s.credentialStorage.Update(ctx, model.Credential{
		ID:       id,
		HashPass: hashPass,
	})
	if err != nil {
		childLog.Error("failed to update password from credential", zap.Error(err))
		return err
	}

	childLog.Info("successfully updated password")

	return nil
}

func (s *Service) UpdateUsername(ctx context.Context, id string, username string) error {
	childLog := s.log.With(
		zap.String("op", "auth.Service.UpdateUsername"),
		zap.String("new_username", username),
	)

	err := s.credentialStorage.Update(ctx, model.Credential{
		ID:       id,
		Username: username,
	})
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return ErrUserAlreadyRegistered
		}
		childLog.Error("failed to update username", zap.Error(err))
		return err
	}

	childLog.Info("successfully updated username")

	return nil
}

func (s *Service) RefreshToken(ctx context.Context, token string) (dto.TokenPair, error) {
	childLog := s.log.With(zap.String("op", "auth.Service.RefreshToken"))

	claims, err := s.verifyAndAddTokenToBlacklist(ctx, token)
	if err != nil {
		childLog.Error("failed to add token to blacklist", zap.Error(err))
		return dto.TokenPair{}, err
	}

	cred, err := s.credentialStorage.GetByID(ctx, claims.Subject)
	if err != nil {
		childLog.Error("failed to get credential by id", zap.Error(err))
		return dto.TokenPair{}, err
	}

	tokens, err := s.newTokenPair(cred, claims.ProfileID)
	if err != nil {
		childLog.Error("failed to create new token pair", zap.Error(err))
		return dto.TokenPair{}, err
	}

	return tokens, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	childLog := s.log.With(zap.String("op", "auth.Service.Logout"))

	_, err := s.verifyAndAddTokenToBlacklist(ctx, token)
	if err != nil {
		childLog.Error("failed to add token to blacklist", zap.Error(err))
		return err
	}

	return nil
}

func (s *Service) verifyAndAddTokenToBlacklist(ctx context.Context, token string) (RefreshClaims, error) {
	claims, err := s.parseRefreshToken(token)
	if err != nil {
		return RefreshClaims{}, err
	}

	exists, err := s.tokenBlacklist.IsExists(ctx, claims.ID)
	if err != nil {
		return RefreshClaims{}, err
	}

	if exists {
		return RefreshClaims{}, ErrTokenBlacklisted
	}

	ttl := claims.ExpiresAt.Time.Sub(time.Now())

	err = s.tokenBlacklist.Add(ctx, claims.ID, ttl)
	if err != nil {
		return RefreshClaims{}, err
	}

	return claims, nil
}
