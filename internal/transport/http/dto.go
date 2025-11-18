// File: internal/transport/http/dto.go
// Purpose: DTOs for auth and user-related requests.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Validation-annotated request bodies for register and login.
package http

// RegisterRequest 注册请求体
// Fields: username(2-64), email(email), password(>=6)
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求体
// Fields: email(email), password(required)
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
