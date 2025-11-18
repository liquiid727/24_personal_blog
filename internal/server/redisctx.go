// File: internal/server/redisctx.go
// Purpose: Helpers to attach Redis client to Gin context.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Provides SetRedis/GetRedis for middleware and handlers.
package server

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetRedis(c *gin.Context, rdb *redis.Client) { c.Set("rdb", rdb) }
func GetRedis(c *gin.Context) *redis.Client {
	v, ok := c.Get("rdb")
	if !ok {
		return nil
	}
	if r, ok := v.(*redis.Client); ok {
		return r
	}
	return nil
}
