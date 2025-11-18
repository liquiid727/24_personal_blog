package server

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetRedis(c *gin.Context, rdb *redis.Client) { c.Set("rdb", rdb) }
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
