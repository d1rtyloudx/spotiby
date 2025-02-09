package image

import (
	"context"
	"encoding/json"
	"github.com/d1rtyloudx/spotiby-pkg/rabbitmq"
	"github.com/d1rtyloudx/spotiby/user-service/internal/config"
	"github.com/d1rtyloudx/spotiby/user-service/internal/domain/model"
	"github.com/d1rtyloudx/spotiby/user-service/internal/dto"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

type imageUploader interface {
	Upload(ctx context.Context, image model.Image) (string, error)
}

type Service struct {
	uploader  imageUploader
	publisher *rabbitmq.Publisher
	cfg       *config.RabbitMQConfig
	log       *zap.Logger
}

func New(uploader imageUploader, publisher *rabbitmq.Publisher, log *zap.Logger, cfg *config.RabbitMQConfig) *Service {
	return &Service{
		uploader:  uploader,
		publisher: publisher,
		log:       log,
		cfg:       cfg,
	}
}

func (s *Service) upload(ctx context.Context, id string, image model.Image, exchange string, routingKey string) error {
	childLog := s.log.With(
		zap.String("op", "image.Service.Upload"),
		zap.String("id", id),
		zap.String("image", image.Name),
		zap.String("content-type", image.ContentType),
		zap.Int64("size", image.Size),
	)

	urlStr, err := s.uploader.Upload(ctx, image)
	if err != nil {
		childLog.Error("failed to upload image", zap.Error(err))
		return err
	}

	req := dto.UpdateAvatarProfileMessage{
		ID:        id,
		AvatarURL: urlStr,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		childLog.Error("failed to marshal request", zap.Error(err))
		return err
	}

	err = s.publisher.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			Body:        reqBytes,
			ContentType: "application/json",
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		childLog.Error("failed to publish image", zap.Error(err))
		return err
	}

	return nil
}

func (s *Service) UploadProfile(ctx context.Context, id string, image model.Image) error {
	return s.upload(
		ctx,
		id,
		image,
		s.cfg.Publishers.ProfileImagePublisher.ExchangeName,
		s.cfg.Publishers.ProfileImagePublisher.RoutingKey,
	)
}

func (s *Service) UploadPlaylist(ctx context.Context, id string, image model.Image) error {
	return s.upload(
		ctx,
		id,
		image,
		s.cfg.Publishers.PlaylistImagePublisher.ExchangeName,
		s.cfg.Publishers.PlaylistImagePublisher.RoutingKey,
	)
}

func (s *Service) UploadTrack(ctx context.Context, id string, image model.Image) error {
	return s.upload(
		ctx,
		id,
		image,
		s.cfg.Publishers.TrackImagePublisher.ExchangeName,
		s.cfg.Publishers.TrackImagePublisher.RoutingKey,
	)
}
