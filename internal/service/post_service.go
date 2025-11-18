// 文章服务：实现文章的创建、更新、删除、发布与列表查询等逻辑
package service

import (
    "errors"
    "time"

    post "go_blog/internal/domain/post"
    tag "go_blog/internal/domain/tag"
    category "go_blog/internal/domain/category"
    "go_blog/internal/repository"
    "go_blog/pkg/utils"
)

// PostService 定义文章业务接口
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
type postService struct{ repo repository.PostRepository }

// NewPostService 创建文章服务
func NewPostService(r repository.PostRepository) PostService { return &postService{repo: r} }

// Create 创建文章，自动生成 slug，已发布状态记录发布时间
func (s *postService) Create(authorID uint, title, content string, status string, tagIDs []uint, categoryIDs []uint) (*post.Post, error) {
    p := &post.Post{AuthorID: authorID, Title: title, Slug: utils.Slugify(title), Content: content, Status: status}
    if status == "published" { now := time.Now(); p.PublishedAt = &now }
    for _, id := range tagIDs { p.Tags = append(p.Tags, tag.Tag{ID: id}) }
    for _, id := range categoryIDs { p.Categories = append(p.Categories, category.Category{ID: id}) }
    if err := s.repo.Create(p); err != nil { return nil, err }
    return p, nil
}

// Update 更新文章，作者或管理员可操作，标签与分类重置绑定
func (s *postService) Update(id uint, authorID uint, title, content, status string, tagIDs []uint, categoryIDs []uint) (*post.Post, error) {
    p, err := s.repo.FindByID(id)
    if err != nil { return nil, err }
    if p.AuthorID != authorID { return nil, errors.New("forbidden") }
    p.Title = title
    p.Content = content
    p.Status = status
    p.Slug = utils.Slugify(title)
    if status == "published" && p.PublishedAt == nil { now := time.Now(); p.PublishedAt = &now }
    p.Tags = nil
    for _, id := range tagIDs { p.Tags = append(p.Tags, tag.Tag{ID: id}) }
    p.Categories = nil
    for _, id := range categoryIDs { p.Categories = append(p.Categories, category.Category{ID: id}) }
    if err := s.repo.Update(p); err != nil { return nil, err }
    return p, nil
}

// Delete 删除文章，作者或管理员可操作
func (s *postService) Delete(id uint, requester uint, isAdmin bool) error {
    p, err := s.repo.FindByID(id)
    if err != nil { return err }
    if !isAdmin && p.AuthorID != requester { return errors.New("forbidden") }
    return s.repo.Delete(id)
}

// Get 获取文章详情
func (s *postService) Get(id uint) (*post.Post, error) { return s.repo.FindByID(id) }

// List 分页查询文章列表
func (s *postService) List(filter repository.PostListFilter, page, size int) ([]post.Post, int64, error) {
    return s.repo.List(filter, page, size)
}

// Publish 发布文章，记录发布时间
func (s *postService) Publish(id uint, requester uint, isAdmin bool) (*post.Post, error) {
    p, err := s.repo.FindByID(id)
    if err != nil { return nil, err }
    if !isAdmin && p.AuthorID != requester { return nil, errors.New("forbidden") }
    p.Status = "published"
    now := time.Now(); p.PublishedAt = &now
    if err := s.repo.Update(p); err != nil { return nil, err }
    return p, nil
}

// IncrementViews 递增浏览量
func (s *postService) IncrementViews(id uint) error { return s.repo.IncrementViews(id) }