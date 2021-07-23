package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/maxeth/go-account-api/model"
)

const (
	TokenRedisSuffix = "refreshtoken"
)

type redisTokenRepository struct {
	Redis *redis.Client
}

func NewTokenRepository(r *redis.Client) model.TokenRepository {
	return &redisTokenRepository{
		Redis: r,
	}
}

func (r *redisTokenRepository) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {

	key := fmt.Sprintf("%s-%s:%s", userID, TokenRedisSuffix, tokenID)

	if err := r.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.Printf("error refreshing token key-value-pair %s:%s in redis repository. error: %v\n", userID, tokenID, err)
		return err
	}

	return nil
}

func (r *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {
	key := fmt.Sprintf("%s-%s:%s", userID, TokenRedisSuffix, tokenID)

	if err := r.Redis.Del(ctx, key).Err(); err != nil {
		log.Printf("error deleting token key-value-pair %s:%s in redis repository. error: %v\n", userID, tokenID, err)
		return err
	}

	return nil
}
