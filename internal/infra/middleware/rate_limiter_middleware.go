package middleware

import (
	"net/http"
)

type RateLimiterMiddleware struct {
	s RateLimiterStrategy
}

func NewRateLimiterMiddleware(strategy RateLimiterStrategy) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{s: strategy}
}

func (md *RateLimiterMiddleware) RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errMsg, statusCode := md.CheckRateLimit(r)
		if errMsg != "" {
			http.Error(w, errMsg, statusCode)
			return
		}

		next.ServeHTTP(w, r)
	})
}