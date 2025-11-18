// File: internal/service/post_service.go
// Purpose: Post business service implementing create/update/delete/publish/list and view counters.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Encapsulates domain logic, slug generation, authorization checks,
//
//	and optional Redis-backed view aggregation for performance.
package service

import (
	"errors"
	"time"

	"context"
	"go_blog/internal/cache"
	category "go_blog/internal/domain/category"
	post "go_blog/internal/domain/post"
	tag "go_blog/internal/domain/tag"
	"go_blog/internal/repository"
	"go_blog/pkg/utils"

	"github.com/redis/go-redis/v9"
)

// PostService 文章业务接口
// Methods:
// - Create: 创建文章，生成 slug，并处理标签与分类绑定
// - Update: 作者或管理员更新，维护发布时间与关联关系
// - Delete: 作者或管理员删除
// - Get/List: 查询详情与分页列表，支持过滤与搜索
// - Publish: 变更状态为 published 并记录时间
// - IncrementViews: 递增浏览量（优先 Redis）
type PostService interface {
	Create(authorID uint, title, content string, status string, tagIDs []uint, categoryIDs []uint) (*post.Post, error)
	Update(id uint, authorID uint, title, content, status string, tagIDs []uint, categoryIDs []uint) (*post.Post, error)
	Delete(id uint, requester uint, isAdmin bool) error
	Get(id uint) (*post.Post, error)
	List(filter repository.PostListFilter, page, size int) ([]post.Post, int64, error)
	Publish(id uint, requester uint, isAdmin bool) (*post.Post, error)
	IncrementViews(id uint) error
}

// postService 文章服务实现
type postService struct {
	repo repository.PostRepository
	rdb  *redis.Client
}

// NewPostService 创建文章服务
func NewPostService(r repository.PostRepository) PostService { return &postService{repo: r} }

func NewPostServiceWithCache(r repository.PostRepository, rdb *redis.Client) PostService {
	return &postService{repo: r, rdb: rdb}
}

// Create 创建文章
// Params: authorID 作者ID，title 标题，content 内容，status 草稿/发布，tagIDs 标签ID集，categoryIDs 分类ID集
// Returns: 新建文章；错误时返回具体原因
// Notes: 当 status=published 时记录发布时间；slug 基于标题生成
func (s *postService) Create(authorID uint, title, content string, status string, tagIDs []uint, categoryIDs []uint) (*post.Post, error) {
	p := &post.Post{AuthorID: authorID, Title: title, Slug: utils.Slugify(title), Content: content, Status: status}
	if status == "published" {
		now := time.Now()
		p.PublishedAt = &now
	}
	for _, id := range tagIDs {
		p.Tags = append(p.Tags, tag.Tag{ID: id})
	}
	for _, id := range categoryIDs {
		p.Categories = append(p.Categories, category.Category{ID: id})
	}
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

// Update 更新文章（作者或管理员）
// Auth: 非管理员必须为作者本人
// Behavior: 更新标题/内容/状态/slug，重置并重建标签与分类关联；首次发布补充发布时间
func (s *postService) Update(id uint, authorID uint, title, content, status string, tagIDs []uint, categoryIDs []uint) (*post.Post, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if p.AuthorID != authorID {
		return nil, errors.New("forbidden")
	}
	p.Title = title
	p.Content = content
	p.Status = status
	p.Slug = utils.Slugify(title)
	if status == "published" && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}
	p.Tags = nil
	for _, id := range tagIDs {
		p.Tags = append(p.Tags, tag.Tag{ID: id})
	}
	p.Categories = nil
	for _, id := range categoryIDs {
		p.Categories = append(p.Categories, category.Category{ID: id})
	}
	if err := s.repo.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

// Delete 删除文章
// Auth: 非管理员必须为作者本人
// Returns: 删除成功或错误
func (s *postService) Delete(id uint, requester uint, isAdmin bool) error {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !isAdmin && p.AuthorID != requester {
		return errors.New("forbidden")
	}
	return s.repo.Delete(id)
}

// Get 获取文章详情
func (s *postService) Get(id uint) (*post.Post, error) { return s.repo.FindByID(id) }

// If Redis is enabled, add cached views to DB value and flush periodically

// List 分页查询文章列表
func (s *postService) List(filter repository.PostListFilter, page, size int) ([]post.Post, int64, error) {
	return s.repo.List(filter, page, size)
}

// Publish 发布文章
// Auth: 非管理员必须为作者本人
// Behavior: 将状态置为 published 并写入当前发布时间
func (s *postService) Publish(id uint, requester uint, isAdmin bool) (*post.Post, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if !isAdmin && p.AuthorID != requester {
		return nil, errors.New("forbidden")
	}
	p.Status = "published"
	now := time.Now()
	p.PublishedAt = &now
	if err := s.repo.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

// IncrementViews 递增浏览量
// Strategy: 若配置了 Redis，使用缓存计数以减少数据库写压力；否则直接更新数据库
func (s *postService) IncrementViews(id uint) error {
	if s.rdb != nil {
		return cache.IncrViews(context.Background(), s.rdb, id)
	}
	return s.repo.IncrementViews(id)
}
