package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type TokenBlacklist struct {
	client *redis.Client
}

func NewTokenBlacklist(client *redis.Client) *TokenBlacklist {
	return &TokenBlacklist{
		client: client,
	}
}

func (t *TokenBlacklist) Add(ctx context.Context, jti string, ttl time.Duration) error {
	const op = "redis.TokenBlacklist.Add"

	err := t.client.Set(ctx, jti, "", ttl).Err()
	if err != nil {
		return fmt.Errorf("%s - t.client.Set: %w", op, err)
	}

	return nil
}

func (t *TokenBlacklist) IsExists(ctx context.Context, jti string) (bool, error) {
	const op = "redis.TokenBlacklist.IsExists"

	err := t.client.Get(ctx, jti).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, fmt.Errorf("%s - t.client.Get: %w", op, err)
	}

	return true, nil
}
