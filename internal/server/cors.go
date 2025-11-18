// File: internal/server/cors.go
// Purpose: CORS middleware handling preflight and allowed origins.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Adds CORS headers; optionally restricts to a single origin.
package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 返回一个处理跨域的中间件
// Params: allowedOrigin 可选的允许来源；为空表示放行请求来源
// Behavior: 设置标准 CORS 头，处理 OPTIONS 预检请求
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
