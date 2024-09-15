package middleware

import (
	"context"
	"net/http"
	"strconv"
)

const (
	rateLimitMsg = "you have reached the maximum number of requests or actions allowed within a certain time frame"
	internalErrMsg = "Internal Server Error"
)

type RateLimiterMiddleware struct {
	s RateLimiterStrategy
}

func NewRateLimiterMiddleware(strategy RateLimiterStrategy) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{s: strategy}
}

func (md *RateLimiterMiddleware) RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		apiKey, clientIP := getCredentials(r)

		limitKey := getLimitKey(apiKey)
		limit, err := md.getLimit(ctx, limitKey)
		if err != nil {
			http.Error(w, internalErrMsg, http.StatusInternalServerError)
			return
		}

		requestsKey := getRequestsKey(apiKey, clientIP)
		errMsg, statusCode := md.getReachedLimit(ctx, requestsKey, limit)
		if errMsg != "" {
			http.Error(w, errMsg, statusCode)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getCredentials(r *http.Request) (string, string) {
	return r.Header.Get("API_KEY"), r.RemoteAddr
}

func getLimitKey(key string) string {
	if key == "" {
		return "limit@default"
	}
	return "limit@" + key
}

func getRequestsKey(key, clientIP string) string {
	if key == "" {
		return "requests@" + clientIP
	}
	return "requests@" + key
}

func (md *RateLimiterMiddleware) getLimit(ctx context.Context, limitKey string) (int64, error) {
	limitStr, err := md.s.Get(ctx, limitKey)
	if err != nil {
		return 0, err
	}
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return limit, nil
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