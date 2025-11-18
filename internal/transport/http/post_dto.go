// File: internal/transport/http/post_dto.go
// Purpose: DTOs for post create/update requests.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Validated request structures for posts including tags and categories.
package http

// PostCreateRequest 创建文章请求体
// Fields: 标题(1-256), 内容(required), 状态(draft/published/private), 标签/分类ID数组
type PostCreateRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=256"`
	Content     string `json:"content" binding:"required"`
	Status      string `json:"status" binding:"required,oneof=draft published private"`
	TagIDs      []uint `json:"tag_ids"`
	CategoryIDs []uint `json:"category_ids"`
}

// PostUpdateRequest 更新文章请求体
// Fields: 同创建；更新时重置关联
type PostUpdateRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=256"`
	Content     string `json:"content" binding:"required"`
	Status      string `json:"status" binding:"required,oneof=draft published private"`
	TagIDs      []uint `json:"tag_ids"`
	CategoryIDs []uint `json:"category_ids"`
}
