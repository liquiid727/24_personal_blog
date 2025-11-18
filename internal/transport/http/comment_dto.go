// File: internal/transport/http/comment_dto.go
// Purpose: DTOs for comment create/update requests.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Request bodies for hierarchical comments.
package http

// CommentCreateRequest 创建评论请求体
// Fields: parent_id(可选), content(required)
type CommentCreateRequest struct {
	ParentID *uint  `json:"parent_id"`
	Content  string `json:"content" binding:"required"`
}

// CommentUpdateRequest 更新评论请求体
// Fields: content(required)
type CommentUpdateRequest struct {
	Content string `json:"content" binding:"required"`
}
