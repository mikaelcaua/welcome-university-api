package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/mikaelcaua/welcome-university-api/internal/httpx"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: map[string][]time.Time{},
		limit:    limit,
		window:   window,
	}
}

func (limiter *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := clientIP(r) + ":" + r.URL.Path
		now := time.Now()
		cutoff := now.Add(-limiter.window)

		limiter.mu.Lock()
		recentRequests := limiter.requests[key]
		keptRequests := recentRequests[:0]
		for _, requestTime := range recentRequests {
			if requestTime.After(cutoff) {
				keptRequests = append(keptRequests, requestTime)
			}
		}
		if len(keptRequests) >= limiter.limit {
			limiter.requests[key] = keptRequests
			limiter.mu.Unlock()
			httpx.RespondError(w, httpx.NewHTTPError(http.StatusTooManyRequests, "Muitas tentativas. Tente novamente em instantes."))
			return
		}
		limiter.requests[key] = append(keptRequests, now)
		limiter.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
