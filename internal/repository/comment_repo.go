// 评论仓储：封装评论的持久化访问与查询
package repository

import (
    "gorm.io/gorm"
    comment "go_blog/internal/domain/comment"
)

// CommentRepository 定义评论相关的数据访问接口
type CommentRepository interface {
    Create(cmt *comment.Comment) error
    Update(cmt *comment.Comment) error
    Delete(id uint) error
    FindByID(id uint) (*comment.Comment, error)
    ListByPost(postID uint) ([]comment.Comment, error)
}

type commentRepository struct{ db *gorm.DB }

// NewCommentRepository 创建评论仓储实现
func NewCommentRepository(db *gorm.DB) CommentRepository { return &commentRepository{db: db} }

// Create 新增评论
func (r *commentRepository) Create(cmt *comment.Comment) error { return r.db.Create(cmt).Error }
// Update 更新评论
func (r *commentRepository) Update(cmt *comment.Comment) error { return r.db.Save(cmt).Error }
// Delete 删除评论
func (r *commentRepository) Delete(id uint) error { return r.db.Delete(&comment.Comment{}, id).Error }

// FindByID 根据主键查询评论
func (r *commentRepository) FindByID(id uint) (*comment.Comment, error) {
    var c comment.Comment
    if err := r.db.First(&c, id).Error; err != nil { return nil, err }
    return &c, nil
}

// ListByPost 查询指定文章的所有评论（按时间/ID排序）
func (r *commentRepository) ListByPost(postID uint) ([]comment.Comment, error) {
    var items []comment.Comment
    if err := r.db.Where("post_id = ?", postID).Order("id ASC").Find(&items).Error; err != nil { return nil, err }
    return items, nil
}