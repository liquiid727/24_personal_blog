// File: internal/repository/post_repo.go
// Purpose: Post repository encapsulating persistence and search.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Provides CRUD, list with driver-specific search (PostgreSQL fulltext or LIKE),
//
//	and atomic views increment.
package repository

import (
	"strings"

	post "go_blog/internal/domain/post"

	"gorm.io/gorm"
)

// PostListFilter 列表查询过滤条件
type PostListFilter struct {
	AuthorID *uint
	Status   string
	Query    string
}

// PostRepository 定义文章相关的数据访问接口
type PostRepository interface {
	Create(p *post.Post) error
	Update(p *post.Post) error
	Delete(id uint) error
	FindByID(id uint) (*post.Post, error)
	List(filter PostListFilter, page, size int) ([]post.Post, int64, error)
	IncrementViews(id uint) error
}

type postRepository struct {
	db     *gorm.DB
	driver string
}

// NewPostRepository 创建文章仓储实现
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db, driver: db.Dialector.Name()}
}

// Create 新增文章
func (r *postRepository) Create(p *post.Post) error { return r.db.Create(p).Error }

// Update 更新文章
func (r *postRepository) Update(p *post.Post) error { return r.db.Save(p).Error }

// Delete 删除文章
func (r *postRepository) Delete(id uint) error { return r.db.Delete(&post.Post{}, id).Error }

// FindByID 根据主键查询文章，预加载标签与分类
func (r *postRepository) FindByID(id uint) (*post.Post, error) {
	var p post.Post
	if err := r.db.Preload("Tags").Preload("Categories").First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// List 按过滤与分页条件查询文章列表
// Filter: AuthorID、Status、Query（PostgreSQL 下使用全文检索，否则 LIKE）
// Pagination: 页码从1开始；返回 items 与 total
func (r *postRepository) List(filter PostListFilter, page, size int) ([]post.Post, int64, error) {
	var items []post.Post
	q := r.db.Model(&post.Post{}).Preload("Tags").Preload("Categories")
	if filter.AuthorID != nil {
		q = q.Where("author_id = ?", *filter.AuthorID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if s := strings.TrimSpace(filter.Query); s != "" {
		if r.driver == "postgres" {
			// Use fulltext search
			q = q.Where("to_tsvector('english', coalesce(title,'') || ' ' || coalesce(content,'')) @@ plainto_tsquery('english', ?)", s)
			q = q.Order("ts_rank(to_tsvector('english', coalesce(title,'') || ' ' || coalesce(content,'')), plainto_tsquery('english', '" + s + "')) DESC, id DESC")
		} else {
			like := "%" + s + "%"
			q = q.Where("title LIKE ? OR content LIKE ?", like, like)
		}
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// IncrementViews 原子递增浏览量
// Implementation: 使用 `UPDATE ... SET views = views + 1` 保证并发安全
func (r *postRepository) IncrementViews(id uint) error {
	return r.db.Model(&post.Post{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + 1")).Error
}
