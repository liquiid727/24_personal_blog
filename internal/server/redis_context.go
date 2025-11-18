// File: internal/server/redis_context.go
// Purpose: Internal helpers to attach Redis client to Gin context.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Lowercase variants setRedis/getRedis used within server package.
package server

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func setRedis(c *gin.Context, rdb *redis.Client) { c.Set("rdb", rdb) }
func getRedis(c *gin.Context) *redis.Client {
	v, ok := c.Get("rdb")
	if !ok {
		return nil
	}
	if r, ok := v.(*redis.Client); ok {
		return r
	}
	return nil
}
