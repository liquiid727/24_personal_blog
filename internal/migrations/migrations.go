package migrations

import (
    "gorm.io/gorm"
)

// ApplyMigrations runs DB-specific migrations beyond AutoMigrate
func ApplyMigrations(db *gorm.DB, driver string) error {
    if driver == "postgres" {
        // Create a GIN index for fulltext search on posts (title + content)
        sql := "CREATE INDEX IF NOT EXISTS idx_posts_search ON posts USING GIN (to_tsvector('english', coalesce(title,'') || ' ' || coalesce(content,'')));"
        if err := db.Exec(sql).Error; err != nil { return err }
    }
    return nil
}