package middleware

import (
	"context"
	"net/http"
	"strconv"
)

const (
	rateLimitMsg = "you have reached the maximum number of requests or actions allowed within a certain time frame"
	invalidKey = "Invalid API Key"
	internalErrMsg = "Internal Server Error"
)


func (md *RateLimiterMiddleware) CheckRateLimit(r *http.Request) (errMsg string, statusCode int) {
	ctx := r.Context()
	apiKey, clientIP := getCredentials(r)

	limitKey := getLimitKey(apiKey)
	limit, errMsg, statusCode := md.getLimit(ctx, limitKey)
	if errMsg != "" {
		return errMsg, statusCode
	}

	requestsKey := getRequestsKey(apiKey, clientIP)
	errMsg, statusCode = md.getReachedLimit(ctx, requestsKey, limit)
	if errMsg != "" {
		return errMsg, statusCode
	}

	return "", 0
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

func (md *RateLimiterMiddleware) getLimit(ctx context.Context, limitKey string) (int64, string, int) {
	limitStr, err := md.s.Get(ctx, limitKey)
	if err != nil {
		return 0, internalErrMsg, http.StatusInternalServerError
	}
	if limitStr == "" {
		return 0, invalidKey, http.StatusUnauthorized
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		return 0, internalErrMsg, http.StatusInternalServerError
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
