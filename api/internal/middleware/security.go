package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
)

// RateLimiter mengelola pembatasan jumlah request per IP dalam satu window waktu (fixed/sliding window).
type RateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	clients map[string]*clientBucket
}

type clientBucket struct {
	count   int
	resetAt time.Time
}

// NewRateLimiter membuat instance RateLimiter (mis. 5 req/menit untuk login, 60 req/menit untuk API umum).
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	if limit <= 0 {
		limit = 60
	}
	if window <= 0 {
		window = time.Minute
	}
	return &RateLimiter{
		limit:   limit,
		window:  window,
		clients: make(map[string]*clientBucket),
	}
}

// Middleware mengembalikan Gin middleware yang membatasi request berdasarkan IP client.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			ip = "unknown"
		}

		now := time.Now()
		rl.mu.Lock()
		bucket, exists := rl.clients[ip]
		if !exists || !now.Before(bucket.resetAt) {
			bucket = &clientBucket{
				count:   0,
				resetAt: now.Add(rl.window),
			}
			rl.clients[ip] = bucket
		}
		bucket.count++
		currentCount := bucket.count
		retryAfterSec := int(time.Until(bucket.resetAt).Seconds())
		if retryAfterSec < 1 {
			retryAfterSec = 1
		}
		rl.mu.Unlock()

		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		remaining := rl.limit - currentCount
		if remaining < 0 {
			remaining = 0
		}
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if currentCount > rl.limit {
			c.Header("Retry-After", strconv.Itoa(retryAfterSec))
			helper.ErrorResponse(c, http.StatusTooManyRequests, fmt.Sprintf("Terlalu banyak permintaan. Silakan coba lagi dalam %d detik.", retryAfterSec))
			c.Abort()
			return
		}

		c.Next()
	}
}

// ErrorHandler menyediakan global error & panic recovery middleware dengan format JSON standar (PRD Backend §8.3).
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				helper.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				c.Abort()
			}
		}()

		c.Next()

		if len(c.Errors) > 0 && !c.Writer.Written() {
			lastErr := c.Errors.Last().Err

			if appErr := helper.GetAppError(lastErr); appErr != nil {
				helper.AppErrorResponse(c, appErr)
				c.Abort()
				return
			}

			switch e := lastErr.(type) {
			case *strconv.NumError:
				helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
			case validator.ValidationErrors:
				details := make(map[string]string, len(e))
				for _, ve := range e {
					details[ve.Field()] = fmt.Sprintf("Validasi gagal pada aturan '%s'", ve.Tag())
				}
				helper.ValidationErrorResponse(c, details)
			default:
				helper.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
			}
			c.Abort()
		}
	}
}
