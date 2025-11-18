package server

import "github.com/gin-gonic/gin"

func New() *gin.Engine {
    r := gin.New()
    r.Use(gin.Recovery())
    return r
}