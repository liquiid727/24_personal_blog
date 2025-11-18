// File: internal/server/middleware.go
// Purpose: Server middlewares including authentication.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Provides JWT auth middleware to validate tokens and inject user context.
package server

import (
	"net/http"
	"strings"

	"go_blog/internal/auth"

	"github.com/gin-gonic/gin"
)

// NewAuthMiddleware 鉴权中间件
// Input: secret JWT HS256 秘钥
// Behavior: 校验 Authorization Bearer Token，失败返回 401；成功注入 uid/role 到上下文
func NewAuthMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(h, "Bearer ")
		claims, err := auth.ParseToken(token, secret)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("uid", claims.UID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
