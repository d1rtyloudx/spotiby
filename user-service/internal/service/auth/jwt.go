package auth

import (
	"fmt"
	"github.com/d1rtyloudx/spotiby/user-service/internal/domain/model"
	"github.com/d1rtyloudx/spotiby/user-service/internal/dto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type AccessClaims struct {
	jwt.RegisteredClaims
	Role        string `json:"role"`
	Email       string `json:"email"`
	ProfileID   string `json:"profile_id"`
	IsConfirmed bool   `json:"is_confirmed"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	ProfileID string `json:"profile_id"`
}

func (s *Service) newAccessToken(credentials model.Credential, profileID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.AccessTTL)),
			Subject:   credentials.ID,
		},
		Role:        credentials.Role,
		Email:       credentials.Email,
		ProfileID:   profileID,
		IsConfirmed: credentials.IsConfirmed,
	})

	tokenString, err := token.SignedString([]byte(s.cfg.AccessSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) newRefreshToken(credentialID string, profileID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.RefreshTTL)),
			Subject:   credentialID,
		},
		ProfileID: profileID,
	})

	tokenString, err := token.SignedString([]byte(s.cfg.RefreshSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) ParseAccessToken(tokenString string) (AccessClaims, error) {
	var claims AccessClaims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.cfg.AccessSecret), nil
	})
	if err != nil {
		return AccessClaims{}, err
	}

	return claims, nil
}

func (s *Service) parseRefreshToken(tokenString string) (RefreshClaims, error) {
	var claims RefreshClaims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.cfg.RefreshSecret), nil
	})
	if err != nil {
		return RefreshClaims{}, err
	}

	return claims, err
}

func (s *Service) newTokenPair(cred model.Credential, profileID string) (dto.TokenPair, error) {
	accessToken, err := s.newAccessToken(cred, profileID)
	if err != nil {
		return dto.TokenPair{}, err
	}

	refreshToken, err := s.newRefreshToken(cred.ID, profileID)
	if err != nil {
		return dto.TokenPair{}, err
	}

	return dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
