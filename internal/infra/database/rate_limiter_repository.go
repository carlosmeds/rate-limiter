package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiterRepository struct {
	RedisClient *redis.Client
}

func NewRateLimiterRepository() *RateLimiterRepository {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &RateLimiterRepository{RedisClient: redisClient}
}

func (r *RateLimiterRepository) Get(ctx context.Context, key string) (string, error) {
	fmt.Println("Getting key", key)
	value, err := r.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Println("Key not found")
		return "", nil
	} else if err != nil {
		fmt.Println("Error", key, err)
		return "", err
	}
	return value, nil
}

func (r *RateLimiterRepository) Save(ctx context.Context, key, value string, ttl int64) error {
	fmt.Println("Saving", value, "for", key)
	err := r.RedisClient.Set(ctx, key, value, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RateLimiterRepository) HasReachedLimit(ctx context.Context, apiKey string, limit int64) (bool, error) {
	count, err := r.RedisClient.Incr(ctx, apiKey).Result()
	if err != nil {
		return false, err
	}

	fmt.Println("Count", count)
	if count == 1 {
		err = r.RedisClient.Expire(ctx, apiKey, 1*time.Second).Err()
		if err != nil {
			return true, err
		}
	}

	if count > limit {
		return true, nil
	}
	return false, nil
}
