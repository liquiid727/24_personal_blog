// 文章仓储：封装文章的持久化访问与查询
package repository

import (
    "strings"

    "gorm.io/gorm"
    post "go_blog/internal/domain/post"
)

// PostListFilter 列表查询过滤条件
type PostListFilter struct {
    AuthorID   *uint
    Status     string
    Query      string
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

type postRepository struct{ db *gorm.DB }

// NewPostRepository 创建文章仓储实现
func NewPostRepository(db *gorm.DB) PostRepository { return &postRepository{db: db} }

// Create 新增文章
func (r *postRepository) Create(p *post.Post) error { return r.db.Create(p).Error }

// Update 更新文章
func (r *postRepository) Update(p *post.Post) error { return r.db.Save(p).Error }

// Delete 删除文章
func (r *postRepository) Delete(id uint) error { return r.db.Delete(&post.Post{}, id).Error }

// FindByID 根据主键查询文章，预加载标签与分类
func (r *postRepository) FindByID(id uint) (*post.Post, error) {
    var p post.Post
    if err := r.db.Preload("Tags").Preload("Categories").First(&p, id).Error; err != nil { return nil, err }
    return &p, nil
}

// List 按过滤与分页条件查询文章列表
func (r *postRepository) List(filter PostListFilter, page, size int) ([]post.Post, int64, error) {
    var items []post.Post
    q := r.db.Model(&post.Post{}).Preload("Tags").Preload("Categories")
    if filter.AuthorID != nil { q = q.Where("author_id = ?", *filter.AuthorID) }
    if filter.Status != "" { q = q.Where("status = ?", filter.Status) }
    if s := strings.TrimSpace(filter.Query); s != "" {
        like := "%" + s + "%"
        q = q.Where("title LIKE ? OR content LIKE ?", like, like)
    }
    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }
    if err := q.Order("id DESC").Offset((page-1)*size).Limit(size).Find(&items).Error; err != nil { return nil, 0, err }
    return items, total, nil
}

// IncrementViews 原子地递增文章浏览量
func (r *postRepository) IncrementViews(id uint) error {
    return r.db.Model(&post.Post{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + 1")).Error
}