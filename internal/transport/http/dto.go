// 认证与用户相关的请求 DTO
package http

// RegisterRequest 注册请求体
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=2,max=64"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求体
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}