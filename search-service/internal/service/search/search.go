package search

import (
	"context"
	"github.com/d1rtyloudx/spotiby-pkg/lib"
	"github.com/d1rtyloudx/spotiby/search-service/internal/model"
	"go.uber.org/zap"
)

type searchStorage interface {
	IndexProfile(ctx context.Context, profile model.Profile) error
	IndexPlaylist(ctx context.Context, playlist model.Playlist) error
	IndexTrack(ctx context.Context, track model.Track) error
	DeleteProfile(ctx context.Context, id string) error
	DeletePlaylist(ctx context.Context, id string) error
	DeleteTrack(ctx context.Context, id string) error
	Search(ctx context.Context, term string, paginationQuery lib.PaginationQuery) (model.SearchResponse, error)
}

type Service struct {
	searchStorage searchStorage
	log           *zap.Logger
}

func New(searchStorage searchStorage, log *zap.Logger) *Service {
	return &Service{
		searchStorage: searchStorage,
		log:           log,
	}
}

func (s *Service) IndexProfile(ctx context.Context, profile model.Profile) error {
	childLog := s.log.With(
		zap.String("op", "service.Search.IndexProfile"),
		zap.Any("data", profile),
	)

	err := s.searchStorage.IndexProfile(ctx, profile)
	if err != nil {
		childLog.Error(
			"failed to index profile",
			zap.Error(err),
		)
		return err
	}

	childLog.Info("successfully indexed profile")

	return nil
}

func (s *Service) IndexPlaylist(ctx context.Context, playlist model.Playlist) error {
	childLog := s.log.With(
		zap.String("op", "service.Search.IndexPlaylist"),
		zap.Any("data", playlist),
	)

	err := s.searchStorage.IndexPlaylist(ctx, playlist)
	if err != nil {
		childLog.Error(
			"failed to index profile",
			zap.Error(err),
		)
		return err
	}

	childLog.Info("successfully indexed playlist")

	return nil
}

func (s *Service) IndexTrack(ctx context.Context, track model.Track) error {
	childLog := s.log.With(
		zap.String("op", "service.Search.IndexTrack"),
		zap.Any("data", track),
	)

	err := s.searchStorage.IndexTrack(ctx, track)
	if err != nil {
		childLog.Error(
			"failed to index profile",
			zap.Error(err),
		)
		return err
	}

	childLog.Info("successfully indexed track")

	return nil
}

func (s *Service) DeleteTrack(ctx context.Context, dto model.DeleteTrack) error {
	childLog := s.log.With(
		zap.String("op", "service.Search.DeleteTrack"),
		zap.String("id", dto.ID),
	)

	err := s.searchStorage.DeleteTrack(ctx, dto.ID)
	if err != nil {
		childLog.Error("failed to delete track", zap.String("id", dto.ID), zap.Error(err))
		return err
	}

	childLog.Info("successfully deleted track")

	return nil
}

func (s *Service) DeletePlaylist(ctx context.Context, dto model.DeletePlaylist) error {
	childLog := s.log.With(
		zap.String("op", "service.Search.DeletePlaylist"),
		zap.String("id", dto.ID),
	)

	err := s.searchStorage.DeletePlaylist(ctx, dto.ID)
	if err != nil {
		childLog.Error("failed to delete playlist", zap.Error(err))
		return err
	}

	childLog.Info("successfully deleted playlist")

	return nil
}

func (s *Service) DeleteProfile(ctx context.Context, dto model.DeleteProfile) error {
	childLog := s.log.With(
		zap.String("op", "service.Search.DeleteProfile"),
		zap.String("id", dto.ID),
	)

	err := s.searchStorage.DeleteProfile(ctx, dto.ID)
	if err != nil {
		childLog.Error("failed to delete profile", zap.Error(err))
		return err
	}

	childLog.Info("successfully deleted profile")

	return nil
}

func (s *Service) Search(ctx context.Context, term string, paginationQuery lib.PaginationQuery) (model.SearchResponse, error) {
	childLog := s.log.With(
		zap.String("op", "service.Search.Search"),
		zap.String("term", term),
	)

	resp, err := s.searchStorage.Search(ctx, term, paginationQuery)
	if err != nil {
		childLog.Error("failed to search profile", zap.Error(err))
		return model.SearchResponse{}, err
	}

	childLog.Info("successfully search profile")

	return resp, err
}
