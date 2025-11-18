// HTTP 响应封装：统一成功与错误返回体结构
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

// OK 成功响应，code=0
func OK(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

// Error 错误响应，status 为 HTTP 状态码，code 为业务错误码
func Error(c *gin.Context, status int, code int, msg string) {
    c.JSON(status, Response{Code: code, Message: msg})
}