// File: internal/server/ratelimit.go
// Purpose: Simple IP-based token bucket rate limiter middleware.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Maintains in-memory buckets per client IP, resets tokens per interval.
package server

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type bucket struct {
	tokens int
	last   time.Time
}

var rl = struct {
	sync.Mutex
	m map[string]*bucket
}{m: map[string]*bucket{}}

// RateLimit 返回限流中间件
// Params: limit 每窗口令牌数；per 窗口时长
// Behavior: 超出令牌时返回 429；窗口轮换时重置令牌
// Note: 进程内存实现，不适用于多实例；可替换为 Redis/漏桶算法
func RateLimit(limit int, per time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		rl.Lock()
		b, ok := rl.m[ip]
		if !ok {
			b = &bucket{tokens: limit, last: time.Now()}
			rl.m[ip] = b
		}
		now := time.Now()
		elapsed := now.Sub(b.last)
		if elapsed > per {
			b.tokens = limit
			b.last = now
		}
		if b.tokens <= 0 {
			rl.Unlock()
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		b.tokens--
		rl.Unlock()
		c.Next()
	}
}
