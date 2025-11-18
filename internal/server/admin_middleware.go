// File: internal/server/admin_middleware.go
// Purpose: Gin middleware enforcing admin-only access control.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Checks request context for role=admin; aborts with 403 if not matched.
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
