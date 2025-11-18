package server

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func CORS(allowedOrigin string) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        if allowedOrigin != "" && origin != allowedOrigin {
            // 简易限制：不匹配则不返回CORS头
            c.Next()
            return
        }
        c.Header("Access-Control-Allow-Origin", origin)
        c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
        c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        c.Next()
    }
}