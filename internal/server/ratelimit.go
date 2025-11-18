package server

import (
    "sync"
    "time"
    "net/http"
    "github.com/gin-gonic/gin"
)

type bucket struct { tokens int; last time.Time }
var rl = struct{ sync.Mutex; m map[string]*bucket }{m: map[string]*bucket{}}

func RateLimit(limit int, per time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        rl.Lock()
        b, ok := rl.m[ip]
        if !ok { b = &bucket{tokens: limit, last: time.Now()}; rl.m[ip] = b }
        now := time.Now()
        elapsed := now.Sub(b.last)
        if elapsed > per { b.tokens = limit; b.last = now }
        if b.tokens <= 0 { rl.Unlock(); c.AbortWithStatus(http.StatusTooManyRequests); return }
        b.tokens--
        rl.Unlock()
        c.Next()
    }
}