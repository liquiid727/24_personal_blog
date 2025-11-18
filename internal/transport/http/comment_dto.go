// 评论相关的请求 DTO
package http

// CommentCreateRequest 创建评论请求体
type CommentCreateRequest struct {
    ParentID *uint  `json:"parent_id"`
    Content  string `json:"content" binding:"required"`
}

// CommentUpdateRequest 更新评论请求体
type CommentUpdateRequest struct {
    Content string `json:"content" binding:"required"`
}