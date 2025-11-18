// 认证处理器：注册、登录与获取当前用户信息
package server

import (
    "net/http"

    "github.com/gin-gonic/gin"
    h "go_blog/internal/transport/http"
    "go_blog/internal/service"
)

// AuthHandler 封装用户认证相关的 HTTP 处理
type AuthHandler struct {
    svc service.UserService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(svc service.UserService) *AuthHandler {
    return &AuthHandler{svc: svc}
}

// Register 注册新用户
func (hnd *AuthHandler) Register(c *gin.Context) {
    var req h.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.Error(c, http.StatusBadRequest, 1, "invalid request")
        return
    }
    u, err := hnd.svc.Register(req.Username, req.Email, req.Password)
    if err != nil {
        h.Error(c, http.StatusBadRequest, 2, "register failed")
        return
    }
    h.OK(c, gin.H{"id": u.ID, "username": u.Username, "email": u.Email})
}

// Login 用户登录并颁发令牌
func (hnd *AuthHandler) Login(c *gin.Context) {
    var req h.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.Error(c, http.StatusBadRequest, 1, "invalid request")
        return
    }
    token, u, err := hnd.svc.Login(req.Email, req.Password)
    if err != nil {
        h.Error(c, http.StatusUnauthorized, 3, "invalid credentials")
        return
    }
    h.OK(c, gin.H{"token": token, "user": gin.H{"id": u.ID, "username": u.Username, "email": u.Email}})
}

// Me 获取当前鉴权用户的基本信息
func (hnd *AuthHandler) Me(c *gin.Context) {
    v, ok := c.Get("uid")
    if !ok {
        h.Error(c, http.StatusUnauthorized, 4, "unauthorized")
        return
    }
    id := uint(v.(uint))
    u, err := hnd.svc.GetByID(id)
    if err != nil {
        h.Error(c, http.StatusNotFound, 5, "not found")
        return
    }
    h.OK(c, gin.H{"id": u.ID, "username": u.Username, "email": u.Email, "role": u.Role})
}