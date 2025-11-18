// 评论服务：实现评论的创建、更新、删除与树形返回逻辑
package service

import (
    "errors"
    comment "go_blog/internal/domain/comment"
    "go_blog/internal/repository"
)

// CommentService 定义评论业务接口
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

// ListTree 构建并返回指定文章的评论树
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