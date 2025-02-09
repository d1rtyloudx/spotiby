package profile

import (
	"context"
	"errors"
	"github.com/d1rtyloudx/spotiby-pkg/lib"
	"go.uber.org/zap"
	"user-service/internal/converter"
	"user-service/internal/domain/model"
	"user-service/internal/dto"
	"user-service/internal/storage"
)

var (
	ErrProfileNotFound = errors.New("profile not found")
)

type profileStorage interface {
	Get(ctx context.Context, pageQuery lib.PaginationQuery) ([]model.Profile, lib.PaginationResponse, error)
	GetByID(ctx context.Context, id string) (model.Profile, error)
	Update(ctx context.Context, profile model.Profile) error
	GetFollows(ctx context.Context, profileID string, pageQuery lib.PaginationQuery) ([]model.Profile, lib.PaginationResponse, error)
	FollowProfile(ctx context.Context, followerID string, followeeID string) error
	UnfollowProfile(ctx context.Context, followerID string, followeeID string) error
}

type Service struct {
	profileStorage profileStorage
	log            *zap.Logger
}

func New(profileStorage profileStorage, log *zap.Logger) *Service {
	return &Service{
		profileStorage: profileStorage,
		log:            log,
	}
}

func (s *Service) Get(ctx context.Context, query lib.PaginationQuery) (dto.PagedProfileResponse, error) {
	profiles, pageResp, err := s.profileStorage.Get(ctx, query)
	if err != nil {
		s.log.Error(
			"failed to get profiles",
			zap.String("op", "profile.Service.Get"),
			zap.Error(err),
		)
		return dto.PagedProfileResponse{}, err
	}

	return dto.PagedProfileResponse{
		Profiles:   converter.ProfileListToProfileDTO(profiles),
		Pagination: pageResp,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (dto.Profile, error) {
	childLog := s.log.With(
		zap.String("op", "profile.Service.GetByID"),
		zap.String("id", id),
	)

	profile, err := s.profileStorage.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return dto.Profile{}, ErrProfileNotFound
		}
		childLog.Error("failed to get profile by id", zap.Error(err))
		return dto.Profile{}, err
	}

	childLog.Info("successfully get profile by id")

	return converter.ProfileToProfileDTO(profile), nil
}

func (s *Service) Update(ctx context.Context, id string, req dto.UpdateProfileRequest) error {
	childLog := s.log.With(
		zap.String("op", "profile.Service.Update"),
		zap.String("id", id),
		zap.Any("data", req),
	)

	err := s.profileStorage.Update(ctx, model.Profile{
		ID:          id,
		DisplayName: req.DisplayName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Description: req.Description,
		AvatarURL:   req.AvatarURL,
	})
	if err != nil {
		childLog.Error("failed to update profile", zap.Error(err))
		return err
	}

	childLog.Info("successfully update profile")

	return nil
}

func (s *Service) GetFollows(ctx context.Context, profileID string, pageQuery lib.PaginationQuery) (dto.PagedProfileResponse, error) {
	profiles, pageResp, err := s.profileStorage.GetFollows(ctx, profileID, pageQuery)
	if err != nil {
		s.log.Error(
			"failed to get profiles",
			zap.String("op", "profileService.GetFollows"),
			zap.Error(err),
		)
		return dto.PagedProfileResponse{}, err
	}

	return dto.PagedProfileResponse{
		Profiles:   converter.ProfileListToProfileDTO(profiles),
		Pagination: pageResp,
	}, nil
}

func (s *Service) FollowProfile(ctx context.Context, followerID string, followeeID string) error {
	childLog := s.log.With(
		zap.String("op", "profile.Service.FollowProfile"),
		zap.String("followerId", followerID),
		zap.String("followeeId", followeeID),
	)

	err := s.profileStorage.FollowProfile(ctx, followerID, followeeID)
	if err != nil {
		childLog.Error("failed to follow profile", zap.Error(err))
		return err
	}

	return nil
}

func (s *Service) UnfollowProfile(ctx context.Context, followerID string, followeeID string) error {
	childLog := s.log.With(
		zap.String("op", "profile.Service.UnfollowProfile"),
		zap.String("followerId", followerID),
		zap.String("followeeId", followeeID),
	)

	err := s.profileStorage.UnfollowProfile(ctx, followerID, followeeID)
	if err != nil {
		childLog.Error("failed to unfollow profile", zap.Error(err))
		return err
	}

	return nil
}
