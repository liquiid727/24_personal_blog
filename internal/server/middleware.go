// 服务中间件：包含认证等通用处理
package server

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "go_blog/internal/auth"
)

// NewAuthMiddleware 校验 Authorization Bearer Token 并将 uid/role 注入上下文
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