// File: internal/config/config.go
// Purpose: Configuration loader for application runtime parameters via environment variables.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Uses Viper to read env with defaults; provides Zap logger helper by env.
package config

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config 定义应用所需的关键配置项
type Config struct {
	AppPort   string
	Env       string
	DBDriver  string
	DBDsn     string
	RedisAddr string
	RedisDB   int
	RedisPass string
	JWTSecret string
	JWTTTL    int
}

// Load 加载配置
// Source: 环境变量（前缀 BLOG_），包含端口、环境、DB、Redis、JWT
// Returns: *Config 或错误
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

// NewLogger 创建日志器
// Params: env 环境字符串（dev/prod）
// Returns: Zap 日志器（开发/生产配置）
func NewLogger(env string) (*zap.Logger, error) {
	if env == "prod" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
