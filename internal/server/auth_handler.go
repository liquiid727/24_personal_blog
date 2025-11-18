// File: internal/server/auth_handler.go
// Purpose: Authentication HTTP handlers for register, login, and current user info.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Server layer handler wiring Gin routes to UserService; validates input,
//
//	handles error responses, and returns normalized JSON payloads.
package server

import (
	"net/http"

	"go_blog/internal/service"
	h "go_blog/internal/transport/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler 封装用户认证相关的 HTTP 处理
type AuthHandler struct {
	svc service.UserService
}

// NewAuthHandler 创建认证处理器
// Params: svc 用户服务实现，用于处理注册、登录与查询
// Returns: *AuthHandler 绑定了依赖的处理器实例
func NewAuthHandler(svc service.UserService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register 注册新用户
// Request: JSON {username,email,password}
// Response: 200 OK -> 基本用户信息；400 -> 请求校验或业务失败
// Edge Cases: 重复邮箱、弱密码策略可在服务层扩展
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
// Request: JSON {email,password}
// Response: 200 OK -> token 与用户信息；401 -> 凭据错误；400 -> 请求格式错误
// Security: 令牌通过服务层签发，遵循配置的 TTL 与秘钥
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
// Auth: 需要上游鉴权中间件注入 uid
// Response: 200 OK -> 用户基本信息与角色；401/404 -> 未鉴权或不存在
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
