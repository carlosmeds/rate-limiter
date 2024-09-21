package database

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiterRepository_Get(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := &RateLimiterRepository{RedisClient: db}
	ctx := context.Background()

	t.Run("key exists", func(t *testing.T) {
		key := "existing_key"
		expectedValue := "some_value"
		mock.ExpectGet(key).SetVal(expectedValue)

		value, err := repo.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	})

	t.Run("key does not exist", func(t *testing.T) {
		key := "non_existing_key"
		mock.ExpectGet(key).RedisNil()

		value, err := repo.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, "", value)
	})

	t.Run("redis error", func(t *testing.T) {
		key := "error_key"
		mock.ExpectGet(key).SetErr(redis.ErrClosed)

		value, err := repo.Get(ctx, key)
		assert.Error(t, err)
		assert.Equal(t, "", value)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRateLimiterRepository_Save(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := &RateLimiterRepository{RedisClient: db}
	ctx := context.Background()

	t.Run("successful save", func(t *testing.T) {
		key := "key_to_save"
		value := "value_to_save"
		ttl := int64(60)

		mock.ExpectSet(key, value, time.Duration(ttl)*time.Second).SetVal("OK")

		err := repo.Save(ctx, key, value, ttl)
		assert.NoError(t, err)
	})

	t.Run("redis error on save", func(t *testing.T) {
		key := "key_to_save"
		value := "value_to_save"
		ttl := int64(60)

		mock.ExpectSet(key, value, time.Duration(ttl)*time.Second).SetErr(redis.ErrClosed)

		err := repo.Save(ctx, key, value, ttl)
		assert.Error(t, err)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRateLimiterRepository_HasReachedLimit(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := &RateLimiterRepository{RedisClient: db}
	ctx := context.Background()

	t.Run("first request within limit", func(t *testing.T) {
		apiKey := "api_key_1"
		limit := int64(5)

		mock.ExpectIncr(apiKey).SetVal(1)
		mock.ExpectExpire(apiKey, 1*time.Second).SetVal(true)

		reachedLimit, err := repo.HasReachedLimit(ctx, apiKey, limit)
		assert.NoError(t, err)
		assert.False(t, reachedLimit)
	})

	t.Run("subsequent request within limit", func(t *testing.T) {
		apiKey := "api_key_2"
		limit := int64(5)

		mock.ExpectIncr(apiKey).SetVal(3)

		reachedLimit, err := repo.HasReachedLimit(ctx, apiKey, limit)
		assert.NoError(t, err)
		assert.False(t, reachedLimit)
	})

	t.Run("request exceeds limit", func(t *testing.T) {
		apiKey := "api_key_3"
		limit := int64(5)

		mock.ExpectIncr(apiKey).SetVal(6)

		reachedLimit, err := repo.HasReachedLimit(ctx, apiKey, limit)
		assert.NoError(t, err)
		assert.True(t, reachedLimit)
	})

	t.Run("redis error on incr", func(t *testing.T) {
		apiKey := "api_key_4"
		limit := int64(5)

		mock.ExpectIncr(apiKey).SetErr(redis.ErrClosed)

		reachedLimit, err := repo.HasReachedLimit(ctx, apiKey, limit)
		assert.Error(t, err)
		assert.False(t, reachedLimit)
	})

	t.Run("redis error on expire", func(t *testing.T) {
		apiKey := "api_key_5"
		limit := int64(5)

		mock.ExpectIncr(apiKey).SetVal(1)
		mock.ExpectExpire(apiKey, 1*time.Second).SetErr(redis.ErrClosed)

		reachedLimit, err := repo.HasReachedLimit(ctx, apiKey, limit)
		assert.Error(t, err)
		assert.True(t, reachedLimit)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

