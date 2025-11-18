// Swagger 路由：提供 OpenAPI 文档与 UI 页面
package server

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// RegisterSwaggerRoutes 注册 OpenAPI JSON 与 Swagger UI 路由
func RegisterSwaggerRoutes(r *gin.Engine, spec []byte) {
    r.GET("/openapi.json", func(c *gin.Context) {
        c.Data(http.StatusOK, "application/json", spec)
    })
    r.GET("/swagger/index.html", func(c *gin.Context) {
        html := `<!doctype html><html><head><link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"></head><body><div id="swagger-ui"></div><script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script><script>window.ui=SwaggerUIBundle({url:'/openapi.json',dom_id:'#swagger-ui'});</script></body></html>`
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
    })
}