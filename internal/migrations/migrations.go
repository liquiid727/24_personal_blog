// File: internal/migrations/migrations.go
// Purpose: Additional driver-specific migrations beyond AutoMigrate.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Adds PostgreSQL GIN index for fulltext search on posts.
package migrations

import (
	"gorm.io/gorm"
)

// ApplyMigrations 执行驱动专属迁移
// Params: db 连接；driver 驱动名
// Behavior: 在 Postgres 下创建全文检索 GIN 索引
func ApplyMigrations(db *gorm.DB, driver string) error {
	if driver == "postgres" {
		// Create a GIN index for fulltext search on posts (title + content)
		sql := "CREATE INDEX IF NOT EXISTS idx_posts_search ON posts USING GIN (to_tsvector('english', coalesce(title,'') || ' ' || coalesce(content,'')));"
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}
