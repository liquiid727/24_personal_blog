// File: internal/server/admin_handler.go
// Purpose: Admin HTTP handler for system statistics.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Requires admin role; aggregates counts of users, posts, comments.
package server

import (
	"go_blog/internal/service"
	h "go_blog/internal/transport/http"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct{ svc service.AdminService }

func NewAdminHandler(s service.AdminService) *AdminHandler { return &AdminHandler{svc: s} }

func (ah *AdminHandler) Stats(c *gin.Context) {
	users, posts, comments, err := ah.svc.Stats()
	if err != nil {
		h.Error(c, http.StatusInternalServerError, 1, "stats failed")
		return
	}
	h.OK(c, gin.H{"users": users, "posts": posts, "comments": comments})
}
