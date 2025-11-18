// 配置模块：负责通过环境变量加载应用运行时所需的参数
package config

import (
    "time"

    "github.com/spf13/viper"
    "go.uber.org/zap"
)

// Config 定义应用所需的关键配置项
type Config struct {
    AppPort    string
    Env        string
    DBDriver   string
    DBDsn      string
    RedisAddr  string
    RedisDB    int
    RedisPass  string
    JWTSecret  string
    JWTTTL     int
}

// Load 从环境变量读取配置，提供默认值，返回 Config
func Load() (*Config, error) {
    v := viper.New()
    v.SetEnvPrefix("BLOG")
    v.AutomaticEnv()

    v.SetDefault("APP_PORT", "8080")
    v.SetDefault("ENV", "dev")
    v.SetDefault("DB_DRIVER", "postgres")
    v.SetDefault("DB_DSN", "postgres://postgres:postgres@localhost:5432/go_blog?sslmode=disable")
    v.SetDefault("REDIS_ADDR", "127.0.0.1:6379")
    v.SetDefault("REDIS_DB", 0)
    v.SetDefault("REDIS_PASS", "")
    v.SetDefault("JWT_SECRET", "change_me")
    v.SetDefault("JWT_TTL", int(time.Hour.Minutes()))

    cfg := &Config{
        AppPort:   v.GetString("APP_PORT"),
        Env:       v.GetString("ENV"),
        DBDriver:  v.GetString("DB_DRIVER"),
        DBDsn:     v.GetString("DB_DSN"),
        RedisAddr: v.GetString("REDIS_ADDR"),
        RedisDB:   v.GetInt("REDIS_DB"),
        RedisPass: v.GetString("REDIS_PASS"),
        JWTSecret: v.GetString("JWT_SECRET"),
        JWTTTL:    v.GetInt("JWT_TTL"),
    }
    return cfg, nil
}

// NewLogger 返回适用于当前环境的 Zap 日志器
func NewLogger(env string) (*zap.Logger, error) {
    if env == "prod" {
        return zap.NewProduction()
    }
    return zap.NewDevelopment()
}