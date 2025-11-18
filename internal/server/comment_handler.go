// File: internal/server/comment_handler.go
// Purpose: Comment HTTP handler for create, update, delete and tree listing.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Handles comment operations with authorization checks and returns
//
//	hierarchical comment trees for posts.
package server

import (
	"net/http"
	"strconv"

	"go_blog/internal/service"
	h "go_blog/internal/transport/http"

	"github.com/gin-gonic/gin"
)

// CommentHandler 封装评论相关的 HTTP 处理
type CommentHandler struct{ svc service.CommentService }

// NewCommentHandler 创建评论处理器
func NewCommentHandler(s service.CommentService) *CommentHandler { return &CommentHandler{svc: s} }

// Create 创建评论，支持多层回复（parent_id 可选）
func (ch *CommentHandler) Create(c *gin.Context) {
	var req h.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, http.StatusBadRequest, 1, "invalid request")
		return
	}
	postID, _ := strconv.Atoi(c.Param("id"))
	uid := c.GetUint("uid")
	cm, err := ch.svc.Create(uint(postID), uid, req.ParentID, req.Content)
	if err != nil {
		h.Error(c, http.StatusBadRequest, 2, "create failed")
		return
	}
	h.OK(c, cm)
}

// Update 更新评论，作者或管理员可操作
func (ch *CommentHandler) Update(c *gin.Context) {
	var req h.CommentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Error(c, http.StatusBadRequest, 1, "invalid request")
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	uid := c.GetUint("uid")
	isAdmin := c.GetString("role") == "admin"
	cm, err := ch.svc.Update(uint(id), uid, req.Content, isAdmin)
	if err != nil {
		h.Error(c, http.StatusForbidden, 3, "update failed")
		return
	}
	h.OK(c, cm)
}

// Delete 删除评论，作者或管理员可操作
func (ch *CommentHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uid := c.GetUint("uid")
	isAdmin := c.GetString("role") == "admin"
	if err := ch.svc.Delete(uint(id), uid, isAdmin); err != nil {
		h.Error(c, http.StatusForbidden, 4, "delete failed")
		return
	}
	h.OK(c, gin.H{"id": id})
}

// ListTree 返回评论树结构
func (ch *CommentHandler) ListTree(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	items, err := ch.svc.ListTree(uint(postID))
	if err != nil {
		h.Error(c, http.StatusBadRequest, 5, "list failed")
		return
	}
	h.OK(c, items)
}
