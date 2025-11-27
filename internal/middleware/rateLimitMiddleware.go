package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter returns a Gin middleware enforcing rate limits per IP.
func RateLimiter(redisURL string, redisPass string,redisdb int,limit int) gin.HandlerFunc {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
		Password: redisPass,
		DB: redisdb,
	})
	ctx := context.Background()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		// Increment request count
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limiter failed"})
			c.Abort()
			return
		}

		// Set expiration on first request
		if count == 1 {
			rdb.Expire(ctx, key, time.Minute)
		}

		if count > int64(limit) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}
