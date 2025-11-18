// File: internal/db/db.go
// Purpose: Database module for opening connections and running unified migrations.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Supports Postgres/MySQL via GORM; auto-migrates domain models and applies driver-specific migrations.
package db

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go_blog/internal/config"
	category "go_blog/internal/domain/category"
	comment "go_blog/internal/domain/comment"
	file "go_blog/internal/domain/file"
	post "go_blog/internal/domain/post"
	tag "go_blog/internal/domain/tag"
	user "go_blog/internal/domain/user"
	"go_blog/internal/migrations"
)

// Open 打开数据库连接
// Params: cfg 配置（驱动与 DSN）；logger 用于记录错误
// Returns: *gorm.DB 或错误；仅支持 postgres/mysql
func Open(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch cfg.DBDriver {
	case "postgres":
		dialector = postgres.Open(cfg.DBDsn)
	case "mysql":
		dialector = mysql.Open(cfg.DBDsn)
	default:
		return nil, errors.New("unsupported db driver")
	}
	return gorm.Open(dialector, &gorm.Config{})
}

// AutoMigrate 自动迁移所有领域模型并执行驱动专属迁移
// Models: User、Post、Tag、Category、PostTag、PostCategory、Comment、File
// Ext: 根据驱动执行全文索引等专属迁移
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&user.User{},
		&post.Post{},
		&tag.Tag{},
		&category.Category{},
		&post.PostTag{},
		&post.PostCategory{},
		&comment.Comment{},
		&file.File{},
	); err != nil {
		return err
	}
	// Run driver-specific migrations
	dial := db.Dialector.Name()
	if err := migrations.ApplyMigrations(db, dial); err != nil {
		return err
	}
	return nil
}
