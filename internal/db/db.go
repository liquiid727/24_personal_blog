// 数据库模块：负责打开连接与统一执行模型迁移
package db

import (
    "errors"

    "go.uber.org/zap"
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "go_blog/internal/config"
    user "go_blog/internal/domain/user"
    post "go_blog/internal/domain/post"
    tag "go_blog/internal/domain/tag"
    category "go_blog/internal/domain/category"
    comment "go_blog/internal/domain/comment"
    file "go_blog/internal/domain/file"
)

// Open 根据配置选择数据库驱动并返回 *gorm.DB
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

// AutoMigrate 执行所有领域模型的自动迁移
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &user.User{},
        &post.Post{},
        &tag.Tag{},
        &category.Category{},
        &post.PostTag{},
        &post.PostCategory{},
        &comment.Comment{},
        &file.File{},
    )
}
