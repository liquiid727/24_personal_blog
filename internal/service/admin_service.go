package service

import (
    "gorm.io/gorm"
    user "go_blog/internal/domain/user"
    post "go_blog/internal/domain/post"
    comment "go_blog/internal/domain/comment"
)

type AdminService interface {
    Stats() (users int64, posts int64, comments int64, err error)
}

type adminService struct{ db *gorm.DB }

func NewAdminService(db *gorm.DB) AdminService { return &adminService{db: db} }

func (s *adminService) Stats() (int64, int64, int64, error) {
    var u, p, c int64
    if err := s.db.Model(&user.User{}).Count(&u).Error; err != nil { return 0, 0, 0, err }
    if err := s.db.Model(&post.Post{}).Count(&p).Error; err != nil { return 0, 0, 0, err }
    if err := s.db.Model(&comment.Comment{}).Count(&c).Error; err != nil { return 0, 0, 0, err }
    return u, p, c, nil
}