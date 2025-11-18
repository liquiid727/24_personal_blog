// File: internal/server/router.go
// Purpose: Minimal router factory for Gin engine.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Provides New() to construct a Gin engine with base recovery middleware.
// File: internal/server/router.go
// Purpose: Minimal router factory for Gin engine.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Provides New() to construct a Gin engine with base recovery middleware.
package server

import "github.com/gin-gonic/gin"

// New 创建基础 Gin 引擎
// Behavior: 初始化并启用 Recovery 中间件
// New 创建基础 Gin 引擎
// Behavior: 初始化并启用 Recovery 中间件
func New() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	return r
}
