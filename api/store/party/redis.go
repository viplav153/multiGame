package party

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"multiGame/api/store"
)

type redisStore struct {
	redis      *redis.Client
	expiration string
}

func New(redis *redis.Client) store.RedisStore {
	return &redisStore{redis: redis}
}

func (r *redisStore) SetKeyValue(ctx context.Context, key string, value []byte) error {
	//  Set a key
	err := r.redis.Set(ctx, key, value, 0).Err()
	if err != nil {
		fmt.Println("Error setting key:", err)
		return err
	}

	return nil
}

func (r *redisStore) GetValue(ctx context.Context, key string) (string, error) {
	// Get the value of a key
	val, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		fmt.Println("Error getting key:", err)
		return "", err
	}

	return val, nil
}

func (r *redisStore) SetKeyValueExpirationSame(ctx context.Context, key string, value []byte) error {
	ttlResult := r.redis.TTL(ctx, "abc")
	if ttlResult.Err() != nil {
		panic(ttlResult.Err())
	}
	expiration, _ := ttlResult.Result()

	//  Set a key
	err := r.redis.Set(ctx, key, value, expiration).Err()
	if err != nil {
		fmt.Println("Error setting key:", err)
		return err
	}

	return nil
}

func (r *redisStore) IsKeyPresent(ctx context.Context, key string) bool {
	existsResult := r.redis.Exists(ctx, key)
	if existsResult.Err() != nil {
		return false
	}

	// If key "abc" exists, existsResult.Val() will be 1, otherwise 0
	if existsResult.Val() == 1 {
		return true
	}

	return false
}
