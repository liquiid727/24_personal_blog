// 文章处理器：提供文章的创建、更新、删除、查询与发布等接口
package server

import (
    "net/http"
    "strconv"
    "context"

    "github.com/gin-gonic/gin"
    h "go_blog/internal/transport/http"
    "go_blog/internal/repository"
    "go_blog/internal/service"
    "go_blog/internal/cache"
)

// PostHandler 封装文章相关的 HTTP 处理
type PostHandler struct{ svc service.PostService }

// NewPostHandler 创建文章处理器
func NewPostHandler(s service.PostService) *PostHandler { return &PostHandler{svc: s} }

// Create 创建文章，需鉴权
func (ph *PostHandler) Create(c *gin.Context) {
    var req h.PostCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil { h.Error(c, http.StatusBadRequest, 1, "invalid request"); return }
    uid := c.GetUint("uid")
    p, err := ph.svc.Create(uid, req.Title, req.Content, req.Status, req.TagIDs, req.CategoryIDs)
    if err != nil { h.Error(c, http.StatusBadRequest, 2, "create failed"); return }
    h.OK(c, p)
}

// Update 更新文章，作者或管理员可操作
func (ph *PostHandler) Update(c *gin.Context) {
    var req h.PostUpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil { h.Error(c, http.StatusBadRequest, 1, "invalid request"); return }
    id, _ := strconv.Atoi(c.Param("id"))
    uid := c.GetUint("uid")
    p, err := ph.svc.Update(uint(id), uid, req.Title, req.Content, req.Status, req.TagIDs, req.CategoryIDs)
    if err != nil { h.Error(c, http.StatusForbidden, 3, "update failed"); return }
    h.OK(c, p)
}

// Delete 删除文章，作者或管理员可操作
func (ph *PostHandler) Delete(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    uid := c.GetUint("uid")
    isAdmin := c.GetString("role") == "admin"
    if err := ph.svc.Delete(uint(id), uid, isAdmin); err != nil { h.Error(c, http.StatusForbidden, 4, "delete failed"); return }
    h.OK(c, gin.H{"id": id})
}

// Get 获取文章详情
func (ph *PostHandler) Get(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    p, err := ph.svc.Get(uint(id))
    if err != nil { h.Error(c, http.StatusNotFound, 5, "not found"); return }
    // 合并 Redis 视图计数
    if rdb := getRedis(c); rdb != nil {
        v, _ := cache.GetViews(context.Background(), rdb, uint(id))
        p.Views += v
    }
    h.OK(c, p)
}

// List 分页查询文章，支持搜索与状态过滤
func (ph *PostHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    f := repository.PostListFilter{Status: c.Query("status"), Query: c.Query("q")}
    items, total, err := ph.svc.List(f, page, size)
    if err != nil { h.Error(c, http.StatusBadRequest, 6, "list failed"); return }
    h.OK(c, gin.H{"items": items, "total": total, "page": page, "page_size": size})
}

// Publish 发布文章，作者或管理员可操作
func (ph *PostHandler) Publish(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    uid := c.GetUint("uid")
    isAdmin := c.GetString("role") == "admin"
    p, err := ph.svc.Publish(uint(id), uid, isAdmin)
    if err != nil { h.Error(c, http.StatusForbidden, 7, "publish failed"); return }
    h.OK(c, p)
}

// IncrViews 递增文章浏览量
func (ph *PostHandler) IncrViews(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := ph.svc.IncrementViews(uint(id)); err != nil { h.Error(c, http.StatusBadRequest, 8, "incr failed"); return }
    h.OK(c, gin.H{"id": id})
}