// File: internal/service/comment_service.go
// Purpose: Comment business service implementing creation, update, deletion and tree building.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Encapsulates authorization checks and builds hierarchical comment trees per post.
package service

import (
    "errors"
    comment "go_blog/internal/domain/comment"
    "go_blog/internal/repository"
)

// CommentService 评论业务接口
// Methods:
// - Create: 新建评论（可选父评论）
// - Update: 作者或管理员更新内容
// - Delete: 作者或管理员删除
// - ListTree: 返回树形结构
type CommentService interface {
    Create(postID, userID uint, parentID *uint, content string) (*comment.Comment, error)
    Update(id uint, userID uint, content string, isAdmin bool) (*comment.Comment, error)
    Delete(id uint, userID uint, isAdmin bool) error
    ListTree(postID uint) ([]CommentNode, error)
}

// commentService 评论服务实现
type commentService struct{ repo repository.CommentRepository }

// NewCommentService 创建评论服务
func NewCommentService(r repository.CommentRepository) CommentService { return &commentService{repo: r} }

// Create 创建评论，默认状态 approved，可扩展审核流程
func (s *commentService) Create(postID, userID uint, parentID *uint, content string) (*comment.Comment, error) {
    c := &comment.Comment{PostID: postID, UserID: userID, ParentID: parentID, Content: content, Status: "approved"}
    if err := s.repo.Create(c); err != nil { return nil, err }
    return c, nil
}

// Update 更新评论，作者或管理员可操作
func (s *commentService) Update(id uint, userID uint, content string, isAdmin bool) (*comment.Comment, error) {
    c, err := s.repo.FindByID(id)
    if err != nil { return nil, err }
    if !isAdmin && c.UserID != userID { return nil, errors.New("forbidden") }
    c.Content = content
    if err := s.repo.Update(c); err != nil { return nil, err }
    return c, nil
}

// Delete 删除评论，作者或管理员可操作
func (s *commentService) Delete(id uint, userID uint, isAdmin bool) error {
    c, err := s.repo.FindByID(id)
    if err != nil { return err }
    if !isAdmin && c.UserID != userID { return errors.New("forbidden") }
    return s.repo.Delete(id)
}

// CommentNode 评论树节点
type CommentNode struct {
    comment.Comment
    Children []CommentNode
}

// ListTree 构建评论树
// Algorithm: 先按 ParentID 聚合，再以根评论为入口递归构建树；复杂度约 O(n)
// Returns: 评论树节点数组；错误时返回存储访问错误
func (s *commentService) ListTree(postID uint) ([]CommentNode, error) {
    items, err := s.repo.ListByPost(postID)
    if err != nil { return nil, err }
    byParent := map[uint][]comment.Comment{}
    roots := []comment.Comment{}
    for _, c := range items {
        if c.ParentID == nil { roots = append(roots, c) } else { byParent[*c.ParentID] = append(byParent[*c.ParentID], c) }
    }
    var build func(c comment.Comment) CommentNode
    build = func(c comment.Comment) CommentNode {
        node := CommentNode{Comment: c}
        children := byParent[c.ID]
        for _, ch := range children { node.Children = append(node.Children, build(ch)) }
        return node
    }
    var res []CommentNode
    for _, r := range roots { res = append(res, build(r)) }
    return res, nil
}