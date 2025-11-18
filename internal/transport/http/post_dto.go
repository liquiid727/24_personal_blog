// 文章相关的请求 DTO
package http

// PostCreateRequest 创建文章请求体
type PostCreateRequest struct {
    Title       string `json:"title" binding:"required,min=1,max=256"`
    Content     string `json:"content" binding:"required"`
    Status      string `json:"status" binding:"required,oneof=draft published private"`
    TagIDs      []uint `json:"tag_ids"`
    CategoryIDs []uint `json:"category_ids"`
}

// PostUpdateRequest 更新文章请求体
type PostUpdateRequest struct {
    Title       string `json:"title" binding:"required,min=1,max=256"`
    Content     string `json:"content" binding:"required"`
    Status      string `json:"status" binding:"required,oneof=draft published private"`
    TagIDs      []uint `json:"tag_ids"`
    CategoryIDs []uint `json:"category_ids"`
}