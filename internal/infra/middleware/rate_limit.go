package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/httpx"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{requests: map[string][]time.Time{}, limit: limit, window: window}
}

func (limiter *RateLimiter) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.ClientIP() + ":" + ctx.FullPath()
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
			httpx.RespondError(ctx, httpx.NewHTTPError(http.StatusTooManyRequests, "Muitas tentativas. Tente novamente em instantes."))
			ctx.Abort()
			return
		}
		limiter.requests[key] = append(keptRequests, now)
		limiter.mu.Unlock()
		ctx.Next()
	}
}
