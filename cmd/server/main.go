// File: cmd/server/main.go
// Purpose: Application entrypoint: load config, init logger/DB, register routes, start HTTP server.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Composes repositories and services, wires Gin handlers and middleware,
//
//	configures Swagger, optional Redis caching, and graceful HTTP server startup.
package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go_blog/internal/cache"
	"go_blog/internal/config"
	"go_blog/internal/db"
	"go_blog/internal/repository"
	"go_blog/internal/server"
	"go_blog/internal/service"

	"github.com/redis/go-redis/v9"
)

// main 初始化并启动服务
// Steps:
// 1) 加载配置与日志
// 2) 打开数据库并迁移
// 3) 组装仓储与服务
// 4) 注册路由与中间件（CORS、鉴权、限流、Swagger）
// 5) 启动 HTTP 服务
func main() {
	// 加载配置（优先环境变量），提供端口、数据库、JWT 秘钥等
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// 初始化结构化日志，根据环境选择开发/生产配置
	logger, err := config.NewLogger(cfg.Env)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// 建立数据库连接（支持 Postgres/MySQL）
	gormDB, err := db.Open(cfg, logger)
	if err != nil {
		logger.Fatal("db open error", zap.Error(err))
	}

	// 执行模型迁移，创建或更新基础表结构
	err = db.AutoMigrate(gormDB)
	if err != nil {
		logger.Fatal("db migrate error", zap.Error(err))
	}

	// 组装用户模块：仓储与服务
	userRepo := repository.NewUserRepository(gormDB)
	userService := service.NewUserService(userRepo, []byte(cfg.JWTSecret), time.Duration(cfg.JWTTTL)*time.Minute)
	// 文章模块：仓储与服务
	postRepo := repository.NewPostRepository(gormDB)
	var postService service.PostService
	var rdb *redis.Client
	if cfg.RedisAddr != "" {
		rdb = cache.New(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)
		postService = service.NewPostServiceWithCache(postRepo, rdb)
	} else {
		postService = service.NewPostService(postRepo)
	}
	// 评论模块：仓储与服务
	commentRepo := repository.NewCommentRepository(gormDB)
	commentService := service.NewCommentService(commentRepo)
	fileRepo := repository.NewFileRepository(gormDB)
	fileService := service.NewFileService(fileRepo)
	adminService := service.NewAdminService(gormDB)

	// 初始化 Gin 与基础中间件
	r := gin.New()
	r.Use(gin.Recovery())
	// CORS 可按需配置具体域名（示例为空表示不启用严格校验）
	r.Use(server.CORS(""))

	// 路由前缀与认证路由
	api := r.Group("/api/v1")
	authH := server.NewAuthHandler(userService)
	api.POST("/auth/register", authH.Register)
	api.POST("/auth/login", authH.Login)

	// JWT 鉴权中间件，解析并注入用户上下文
	authMW := server.NewAuthMiddleware([]byte(cfg.JWTSecret))
	apiAuth := api.Group("")
	apiAuth.Use(authMW)
	apiAuth.GET("/auth/me", authH.Me)

	// 文章路由（部分需鉴权），支持 CRUD、发布与浏览量递增
	postH := server.NewPostHandler(postService)
	// 登录、创建评论、发布等敏感操作做速率限制
	apiAuth.POST("/posts", server.RateLimit(30, time.Minute), postH.Create)
	apiAuth.PUT("/posts/:id", postH.Update)
	apiAuth.DELETE("/posts/:id", postH.Delete)
	api.GET("/posts", postH.List)
	api.GET("/posts/:id", postH.Get)
	apiAuth.POST("/posts/:id/publish", postH.Publish)
	api.POST("/posts/:id/views/incr", postH.IncrViews)

	// 评论路由（创建/更新/删除需鉴权），列表提供树形结构
	commentH := server.NewCommentHandler(commentService)
	apiAuth.POST("/posts/:id/comments", server.RateLimit(60, time.Minute), commentH.Create)
	apiAuth.PUT("/comments/:id", server.RateLimit(60, time.Minute), commentH.Update)
	apiAuth.DELETE("/comments/:id", server.RateLimit(60, time.Minute), commentH.Delete)
	api.GET("/posts/:id/comments", commentH.ListTree)

	// 文件上传与管理
	r.Static("/uploads", "./uploads")
	fileH := server.NewFileHandler(fileService)
	apiAuth.POST("/files", fileH.Upload)
	api.GET("/files/:id", fileH.Get)
	apiAuth.DELETE("/files/:id", fileH.Delete)

	// 管理员功能与统计
	adminH := server.NewAdminHandler(adminService)
	admin := apiAuth.Group("/admin")
	admin.Use(server.AdminOnly())
	admin.GET("/stats", adminH.Stats)

	// 注册 Swagger UI 与 OpenAPI JSON（开发环境使用）
	spec, _ := os.ReadFile("docs/openapi.json")
	server.RegisterSwaggerRoutes(r, spec)

	// 启动 HTTP 服务器
	if rdb != nil {
		r.Use(func(c *gin.Context) { server.SetRedis(c, rdb); c.Next() })
	}
	s := &http.Server{Addr: ":" + cfg.AppPort, Handler: r}
	if err := s.ListenAndServe(); err != nil {
		logger.Fatal("server error", zap.Error(err))
	}
}
