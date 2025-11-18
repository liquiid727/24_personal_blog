// 文章领域模型：管理文章内容、状态、标签与分类
package post

import (
    "time"
    tag "go_blog/internal/domain/tag"
    category "go_blog/internal/domain/category"
)

// Post 表示一篇博客文章
type Post struct {
    ID          uint      `gorm:"primaryKey"`
    AuthorID    uint      `gorm:"index"`                 // 作者用户 ID
    Title       string    `gorm:"size:256"`              // 标题
    Slug        string    `gorm:"uniqueIndex;size:256"`  // 唯一短标识
    Content     string    `gorm:"type:text"`             // 正文内容
    Status      string    `gorm:"size:16"`               // 状态（draft/published/private）
    Views       int64     `gorm:"index"`                 // 浏览量
    PublishedAt *time.Time
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Tags        []tag.Tag          `gorm:"many2many:post_tags"`         // 多对多标签
    Categories  []category.Category `gorm:"many2many:post_categories"`  // 多对多分类
}

// PostTag 文章与标签的关联表（多对多）
type PostTag struct {
    PostID uint `gorm:"primaryKey"`
    TagID  uint `gorm:"primaryKey"`
}

// PostCategory 文章与分类的关联表（多对多）
type PostCategory struct {
    PostID     uint `gorm:"primaryKey"`
    CategoryID uint `gorm:"primaryKey"`
}