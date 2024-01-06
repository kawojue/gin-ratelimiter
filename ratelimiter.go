// package rate

// import (
// 	"math/rand"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"golang.org/x/time/rate"
// )

// // LimiterConfig holds the configuration parameters for the rate limiter.
// type LimiterConfig struct {
// 	MaxAttempts int    // Maximum number of attempts allowed within the specified time window.
// 	Message     string // Error message to be sent when the rate limit is exceeded.
// 	TimerArray  []int  // Array of time intervals in seconds for adjusting burst size periodically.
// }

// // Creates and returns a new rate limiter based on the provided configuration.
// // It also starts a goroutine to periodically adjust the burst size according to the random timer.
// func CreateLimiter(config *LimiterConfig) *rate.Limiter {
// 	randomTimer := time.Duration(config.TimerArray[rand.Intn(len(config.TimerArray))]) * time.Second
// 	limiter := rate.NewLimiter(rate.Limit(config.MaxAttempts), int(config.MaxAttempts))

// 	limiter.SetBurst(int(config.MaxAttempts))

// 	go func() {
// 		for {
// 			time.Sleep(randomTimer)
// 			limiter.SetBurst(int(config.MaxAttempts))
// 		}
// 	}()

// 	return limiter
// }

// //	Returns a Gin middleware that checks if a request is allowed by the rate limiter.
// //
// // If not allowed, it responds with a HTTP 429 Too Many Requests status and the specified error message.
// func RateLimiter(limiter *rate.Limiter, config *LimiterConfig) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		if !limiter.Allow() {
// 			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": config.Message})
// 		}

// 		ctx.Next()
// 	}
// }

// Package rate provides rate limiting middleware for the Gin framework using golang.org/x/time/rate.
package rate

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// LimiterConfig holds the configuration parameters for the rate limiter.
type LimiterConfig struct {
	MaxAttempts int    // Maximum number of attempts allowed within the specified time window.
	Message     string // Error message to be sent when the rate limit is exceeded.
	TimerArray  []int  // Array of time intervals in seconds for adjusting burst size periodically.
}

// Creates and returns a new rate limiter based on the provided configuration.
// It also starts a goroutine to periodically adjust the burst size according to the random timer.
func CreateLimiter(config *LimiterConfig) *rate.Limiter {
	limiter := rate.NewLimiter(rate.Limit(config.MaxAttempts), int(config.MaxAttempts))

	go func() {
		ticker := time.NewTicker(time.Second) // Adjust burst size every second
		defer ticker.Stop()

		for {
			<-ticker.C
			limiter.SetBurst(int(config.MaxAttempts))
		}
	}()

	return limiter
}

//	Returns a Gin middleware that checks if a request is allowed by the rate limiter.
//
// If not allowed, it responds with a HTTP 429 Too Many Requests status and the specified error message.
func RateLimiter(limiter *rate.Limiter, config *LimiterConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !limiter.Allow() {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": config.Message})
		}

		ctx.Next()
	}
}
