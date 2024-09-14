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
		ctx := r.Context()
		key := r.Header.Get("API_KEY")

		limit := int64(3)
		reached_limit, err := md.s.HasReachedLimit(ctx, key, limit)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if reached_limit {
			http.Error(w, "TooManyRequests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
