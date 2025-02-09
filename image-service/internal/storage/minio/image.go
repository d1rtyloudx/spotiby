package minio

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"image-service/internal/domain/model"
)

type ImageStorage struct {
	client   *minio.Client
	endpoint string
}

func NewImageStorage(client *minio.Client, endpoint string) *ImageStorage {
	return &ImageStorage{
		client:   client,
		endpoint: endpoint,
	}
}

func (s *ImageStorage) Upload(ctx context.Context, image model.Image) (string, error) {
	const op = "minio.ImageStorage.Upload"

	opts := minio.PutObjectOptions{
		ContentType: image.ContentType,
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
	}

	info, err := s.client.PutObject(
		ctx,
		image.BucketName,
		s.generateFilename(image.Name),
		image.File,
		image.Size,
		opts,
	)
	if err != nil {
		return "", fmt.Errorf("%s - s.client.PutObject: %w", op, err)
	}

	urlString := s.generateURLStr(info.Bucket, info.Key)

	return urlString, nil
}

func (s *ImageStorage) generateURLStr(bucket string, key string) string {
	return fmt.Sprintf("%s/%s/%s", s.endpoint, bucket, key)
}

func (s *ImageStorage) generateFilename(filename string) string {
	return fmt.Sprintf("%s-%s", uuid.New().String(), filename)
}
