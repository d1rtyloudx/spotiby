package app

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

//add methods to spotiby-pkg

func (a *App) mustInitBuckets(ctx context.Context) {
	if err := a.initBuckets(ctx); err != nil {
		panic(err)
	}
}

func (a *App) initBuckets(ctx context.Context) error {
	buckets := []string{
		a.cfg.Minio.Buckets.PlaylistBucket,
		a.cfg.Minio.Buckets.ProfileBucket,
		a.cfg.Minio.Buckets.TrackBucket,
	}

	for _, bucket := range buckets {
		if err := a.createAndSetPolicy(ctx, bucket); err != nil {
			a.log.Error(
				"minio bucket creation failed",
				zap.String("op", "app.InitBuckets"),
				zap.String("bucket", bucket),
				zap.Error(err),
			)
			return err
		}
	}

	a.log.Info("successfully init minio buckets")
	return nil
}

func (a *App) createAndSetPolicy(ctx context.Context, bucketName string) error {
	const op = "app.createAndSetPolicy"

	exists, err := a.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("%s - a.client.BucketExists: %w", op, err)
	}

	if !exists {
		if err := a.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("%s - a.client.MakeBucket, bucket: %s, %w", op, bucketName, err)
		}

		policy := fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": "*",
                "Action": ["s3:GetObject"],
                "Resource": ["arn:aws:s3:::%s/*"]
            }
        ]
    }`, bucketName)

		if err := a.client.SetBucketPolicy(ctx, bucketName, policy); err != nil {
			return fmt.Errorf("%s - a.client.SetBucketPolicy, bucket: %s: %w", op, bucketName, err)
		}
	}

	return nil
}
