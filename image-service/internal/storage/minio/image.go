package minio

import (
	"context"
	"fmt"
	"github.com/d1rtyloudx/spotiby/user-service/internal/domain/model"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
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

func (s *ImageStorage) Get(ctx context.Context, bucketName string, fileName string) (*minio.Object, error) {
	const op = "minio.ImageStorage.Get"

	obj, err := s.client.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Printf("GetObject error: %v\n", err)
		return nil, fmt.Errorf("s.client.GetObject - %s: %w", op, err)
	}

	return obj, nil
}

func (s *ImageStorage) Upload(ctx context.Context, image model.Image) (string, error) {
	const op = "minio.ImageStorage.Upload"

	opts := minio.PutObjectOptions{
		ContentType: image.ContentType,
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
	}

	fileName := s.generateFilename(image.Name)

	_, err := s.client.PutObject(
		ctx,
		image.BucketName,
		fileName,
		image.File,
		image.Size,
		opts,
	)
	if err != nil {
		return "", fmt.Errorf("%s - s.client.PutObject: %w", op, err)
	}

	return fileName, nil
}

func (s *ImageStorage) generateFilename(filename string) string {
	return fmt.Sprintf("%s-%s", uuid.New().String(), filename)
}
