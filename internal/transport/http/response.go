// File: internal/transport/http/response.go
// Purpose: HTTP response helpers to standardize success and error payloads.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Defines Response struct and OK/Error helpers for consistent API outputs.
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OK 成功响应
// Params: data 任意数据；返回 HTTP 200，业务码固定为 0
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

// Error 错误响应
// Params: status HTTP 状态码；code 业务错误码；msg 人类可读消息
func Error(c *gin.Context, status int, code int, msg string) {
	c.JSON(status, Response{Code: code, Message: msg})
}
