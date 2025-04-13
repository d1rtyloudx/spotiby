package image

import (
	"context"
	"encoding/json"
	"github.com/d1rtyloudx/spotiby-pkg/rabbitmq"
	"github.com/d1rtyloudx/spotiby/user-service/internal/config"
	"github.com/d1rtyloudx/spotiby/user-service/internal/domain/model"
	"github.com/d1rtyloudx/spotiby/user-service/internal/dto"
	"github.com/minio/minio-go/v7"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

type imageUploader interface {
	Upload(ctx context.Context, image model.Image) (string, error)
}

type imageProvider interface {
	Get(ctx context.Context, bucketName, fileName string) (*minio.Object, error)
}

type Service struct {
	uploader  imageUploader
	provider  imageProvider
	publisher *rabbitmq.Publisher
	cfg       *config.RabbitMQConfig
	log       *zap.Logger
}

func New(
	uploader imageUploader,
	provider imageProvider,
	publisher *rabbitmq.Publisher,
	log *zap.Logger,
	cfg *config.RabbitMQConfig,
) *Service {
	return &Service{
		uploader:  uploader,
		provider:  provider,
		publisher: publisher,
		log:       log,
		cfg:       cfg,
	}
}

func (s *Service) get(ctx context.Context, bucketName, fileName string) (*minio.Object, error) {
	childLog := s.log.With(
		zap.String("op", "image.Service.get"),
		zap.String("bucket_name", bucketName),
		zap.String("file_name", fileName),
	)

	obj, err := s.provider.Get(ctx, bucketName, fileName)
	if err != nil {
		childLog.Error("failed to get image", zap.Error(err))
		return nil, err
	}

	return obj, err
}

func (s *Service) upload(ctx context.Context, id string, image model.Image, exchange string, routingKey string) (string, error) {
	childLog := s.log.With(
		zap.String("op", "image.Service.upload"),
		zap.String("id", id),
		zap.String("image", image.Name),
		zap.String("content-type", image.ContentType),
		zap.Int64("size", image.Size),
	)

	urlStr, err := s.uploader.Upload(ctx, image)
	if err != nil {
		childLog.Error("failed to upload image", zap.Error(err))
		return "", err
	}

	req := dto.UpdateAvatarProfileMessage{
		ID:        id,
		AvatarURL: urlStr,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		childLog.Error("failed to marshal request", zap.Error(err))
		return "", err
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
		return "", err
	}

	return urlStr, nil
}

func (s *Service) GetProfile(ctx context.Context, buketName, fileName string) (*minio.Object, error) {
	return s.get(ctx, buketName, fileName)
}

func (s *Service) GetTrack(ctx context.Context, buketName, fileName string) (*minio.Object, error) {
	return s.get(ctx, buketName, fileName)
}

func (s *Service) GetPlaylist(ctx context.Context, buketName, fileName string) (*minio.Object, error) {
	return s.get(ctx, buketName, fileName)
}

func (s *Service) UploadProfile(ctx context.Context, id string, image model.Image) (string, error) {
	return s.upload(
		ctx,
		id,
		image,
		s.cfg.Publishers.ProfileImagePublisher.ExchangeName,
		s.cfg.Publishers.ProfileImagePublisher.RoutingKey,
	)
}

func (s *Service) UploadPlaylist(ctx context.Context, id string, image model.Image) (string, error) {
	return s.upload(
		ctx,
		id,
		image,
		s.cfg.Publishers.PlaylistImagePublisher.ExchangeName,
		s.cfg.Publishers.PlaylistImagePublisher.RoutingKey,
	)
}

func (s *Service) UploadTrack(ctx context.Context, id string, image model.Image) (string, error) {
	return s.upload(
		ctx,
		id,
		image,
		s.cfg.Publishers.TrackImagePublisher.ExchangeName,
		s.cfg.Publishers.TrackImagePublisher.RoutingKey,
	)
}
