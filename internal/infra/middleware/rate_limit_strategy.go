package middleware

import (
	"context"

	"github.com/carlosmeds/rate-limiter/internal/infra/database"
)

type RateLimiterStrategy interface {
	HasReachedLimit(ctx context.Context, apiKey string, limit int64) (bool, error)
	Get(ctx context.Context, key string) (string, error)
	Save(ctx context.Context, key, value string, ttl int64) error
}

func NewRateLimiterStrategy() RateLimiterStrategy {
	return database.NewRateLimiterRepository()
}
