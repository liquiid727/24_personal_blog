package cache

import (
    "context"
    "strconv"
    "github.com/redis/go-redis/v9"
)

func New(addr, pass string, db int) *redis.Client {
    return redis.NewClient(&redis.Options{Addr: addr, Password: pass, DB: db})
}

func IncrViews(ctx context.Context, rdb *redis.Client, postID uint) error {
    key := viewsKey(postID)
    return rdb.Incr(ctx, key).Err()
}

func GetViews(ctx context.Context, rdb *redis.Client, postID uint) (int64, error) {
    v, err := rdb.Get(ctx, viewsKey(postID)).Int64()
    if err == redis.Nil { return 0, nil }
    return v, err
}

func ResetViews(ctx context.Context, rdb *redis.Client, postID uint) error {
    key := viewsKey(postID)
    _, err := rdb.Del(ctx, key).Result()
    return err
}

func viewsKey(id uint) string { return "post:views:" + strconv.FormatUint(uint64(id), 10) }