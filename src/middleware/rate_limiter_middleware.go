package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

// RateLimiter 用于存储每个 IP 地址的请求记录
type RateLimiter struct {
	mu         sync.Mutex
	requests   map[string][]time.Time
	limit      int
	windowSize time.Duration
}

// NewRateLimiter 创建一个新的 RateLimiter 实例
func NewRateLimiter(limit int, windowSize time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:   make(map[string][]time.Time),
		limit:      limit,
		windowSize: windowSize,
	}
}

// Allow 判断某个 IP 地址是否允许请求
func (r *RateLimiter) Allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	requests, exists := r.requests[ip]

	if !exists {
		r.requests[ip] = []time.Time{now}
		return true
	}

	// 移除过期的请求
	var recentRequests []time.Time
	for _, reqTime := range requests {
		if now.Sub(reqTime) <= r.windowSize {
			recentRequests = append(recentRequests, reqTime)
		}
	}

	// 检查请求数是否超过限制
	if len(recentRequests) >= r.limit {
		return false
	}

	// 记录当前请求
	r.requests[ip] = append(recentRequests, now)
	return true
}

// RateLimiterMiddleware 用于限制请求频率
func RateLimiterMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "请求过于频繁，请稍后再试"})
			c.Abort()
			return
		}
		c.Next()
	}
}
