package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/carlosmeds/rate-limiter/configs"
)

const (
	rateLimitMsg   = "you have reached the maximum number of requests or actions allowed within a certain time frame"
	invalidKey     = "Invalid API Key"
	internalErrMsg = "Internal Server Error"
)

func (md *RateLimiterMiddleware) CheckRateLimit(r *http.Request) (errMsg string, statusCode int) {
	ctx := r.Context()
	apiKey, clientIP := getCredentials(r)

	blackListKey := getBlackListKey(apiKey, clientIP)
	errMsg, statusCode = md.isBlackListed(ctx, blackListKey)
	if errMsg != "" {
		return errMsg, statusCode
	}

	config := configs.GetConfig()
	limit, errMsg, statusCode := md.getLimit(apiKey, config)
	if errMsg != "" {
		return errMsg, statusCode
	}

	requestsKey := getRequestsKey(apiKey, clientIP)
	errMsg, statusCode = md.getReachedLimit(ctx, requestsKey, limit)
	if errMsg == rateLimitMsg {
		md.AddToBlackList(ctx, blackListKey, config)
		return errMsg, statusCode
	}
	if errMsg != "" {
		return errMsg, statusCode
	}

	return "", 0
}

func getCredentials(r *http.Request) (string, string) {
	return r.Header.Get("API_KEY"), r.RemoteAddr
}

func getRequestsKey(key, clientIP string) string {
	if key == "" {
		return "requests@" + clientIP
	}
	return "requests@" + key
}

func getBlackListKey(key, clientIP string) string {
	if key == "" {
		return "blacklist@" + clientIP
	}
	return "blacklist@" + key
}

func (md *RateLimiterMiddleware) isBlackListed(ctx context.Context, key string) (string, int) {
	blackListed, err := md.s.Get(ctx, key)
	if err != nil {
		return internalErrMsg, http.StatusInternalServerError
	}
	if blackListed != "" {
		return rateLimitMsg, http.StatusTooManyRequests
	}

	return "", 0
}

func (md *RateLimiterMiddleware) getLimit(apiKey string, config *configs.Config) (int64, string, int) {
	if apiKey == "" {
		return config.DefaultLimit, "", 0
	}

	limit, exists := config.ApiKeyLimits[apiKey]
	if !exists {
		return 0, invalidKey, http.StatusUnauthorized
	}

	return limit, "", 0
}

func (md *RateLimiterMiddleware) getReachedLimit(ctx context.Context, key string, limit int64) (string, int) {
	reachedLimit, err := md.s.HasReachedLimit(ctx, key, limit)
	if err != nil {
		return internalErrMsg, http.StatusInternalServerError
	}
	if reachedLimit {
		return rateLimitMsg, http.StatusTooManyRequests
	}
	return "", 0
}

func (md *RateLimiterMiddleware) AddToBlackList(ctx context.Context, key string,  config *configs.Config) error {
	err := md.s.Save(ctx, key, "Too many requests", config.BlockedTime)
	if err != nil {
		fmt.Println("Error adding to blacklist", err)
		return err
	}

	return nil
}
