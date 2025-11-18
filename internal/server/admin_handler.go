package server

import (
    "net/http"
    "github.com/gin-gonic/gin"
    h "go_blog/internal/transport/http"
    "go_blog/internal/service"
)

type AdminHandler struct{ svc service.AdminService }

func NewAdminHandler(s service.AdminService) *AdminHandler { return &AdminHandler{svc: s} }

func (ah *AdminHandler) Stats(c *gin.Context) {
    users, posts, comments, err := ah.svc.Stats()
    if err != nil { h.Error(c, http.StatusInternalServerError, 1, "stats failed"); return }
    h.OK(c, gin.H{"users": users, "posts": posts, "comments": comments})
}