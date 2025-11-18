// File: internal/cache/redis.go
// Purpose: Redis client factory and helpers for post views counter.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Provides New client and functions to increment/get/reset views with namespaced keys.
package cache

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// New 创建 Redis 客户端
// Params: addr 地址；pass 密码；db DB索引
func New(addr, pass string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: addr, Password: pass, DB: db})
}

// IncrViews 递增文章浏览量缓存
// Key: post:views:{id}
func IncrViews(ctx context.Context, rdb *redis.Client, postID uint) error {
	key := viewsKey(postID)
	return rdb.Incr(ctx, key).Err()
}

// GetViews 获取文章浏览量缓存
// Returns: 若键不存在返回 0,nil
func GetViews(ctx context.Context, rdb *redis.Client, postID uint) (int64, error) {
	v, err := rdb.Get(ctx, viewsKey(postID)).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return v, err
}

// ResetViews 重置浏览量缓存
func ResetViews(ctx context.Context, rdb *redis.Client, postID uint) error {
	key := viewsKey(postID)
	_, err := rdb.Del(ctx, key).Result()
	return err
}

func viewsKey(id uint) string { return "post:views:" + strconv.FormatUint(uint64(id), 10) }
