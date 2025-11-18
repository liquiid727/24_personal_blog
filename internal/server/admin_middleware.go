package server

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.GetString("role") != "admin" {
            c.AbortWithStatus(http.StatusForbidden)
            return
        }
        c.Next()
    }
}