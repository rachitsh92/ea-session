package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenBucket represents a token bucket for rate limiting
type TokenBucket struct {
	capacity      int //maximum number of requests a bucket can hold
	tokens        int // current number of tokens
	refillRate    int // number of tokens to add at each refill interval
	refillInterval time.Duration
	mu            sync.Mutex // ensure thread-safety on bucket
}

// NewTokenBucket creates a new TokenBucket
func NewTokenBucket(capacity int, refillRate int, refillInterval time.Duration) *TokenBucket {
	tb := &TokenBucket{
		capacity:      capacity,
		tokens:        capacity,
		refillRate:    refillRate,
		refillInterval: refillInterval,
	}

	go tb.refill() // Starts a goroutine to refill the bucket periodically.
	return tb
}

// refill adds tokens to the bucket at the specified interval
func (tb *TokenBucket) refill() {
	ticker := time.NewTicker(tb.refillInterval) // initiating periodic execution at specified interval
	for range ticker.C { // way of listening to a tick 
		tb.mu.Lock() // ensure thread-safety incase of multiple goroutines
		if tb.tokens < tb.capacity {
			tb.tokens += tb.refillRate
			if tb.tokens > tb.capacity {
				tb.tokens = tb.capacity
			}
		}
		tb.mu.Unlock()
	}
}

// Allow checks if a token is available and consumes one if available
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(tb *TokenBucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		/* 
		User level rate limiting
		user := c.ClientIP() Example: using IP address as the user identifier
		bucket := rl.GetBucket(user)
		*/
		
		/* 
		Similar allow can be used to:
		1. ration compute to individual tier users.
		2. This can even be used for internal ration where you would not want your api 
			to hit the thrid party api's more frequently than you anticipate. 
			(ex: third party AI/ML costly api)
		*/
		if !tb.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	bucket := NewTokenBucket(5, 1, time.Second)

	r := gin.Default()

	r.Use(RateLimitMiddleware(bucket))

	r.GET("/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Request successful",
		})
	})

	r.Run(":8080")
}
