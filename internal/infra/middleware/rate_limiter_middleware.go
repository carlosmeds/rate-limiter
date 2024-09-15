package middleware

import (
	"net/http"
)

const (
	rate_limit_msg = "you have reached the maximum number of requests or actions allowed within a certain time frame"
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
		key := r.Header.Get("API_KEY")

		limit := int64(3)
		reached_limit, err := md.s.HasReachedLimit(ctx, key, limit)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if reached_limit {
			http.Error(w, rate_limit_msg, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
