// File: internal/service/user_service.go
// Purpose: User business service providing registration, login and retrieval.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Validates credentials, hashes passwords, issues JWT tokens per config.
package service

import (
	"errors"
	"time"

	"go_blog/internal/auth"
	"go_blog/internal/domain/user"
	"go_blog/internal/repository"
)

// UserService 用户业务接口
// Methods:
// - Register: 创建用户并校验输入
// - Login: 验证密码并签发 JWT
// - GetByID: 按主键获取用户
type UserService interface {
	Register(username, email, password string) (*user.User, error)
	Login(email, password string) (string, *user.User, error)
	GetByID(id uint) (*user.User, error)
}

// userService 用户服务实现
type userService struct {
	repo      repository.UserRepository
	jwtSecret []byte
	jwtTTL    time.Duration
}

// NewUserService 创建用户服务
// Params: repo 用户仓储；jwtSecret JWT秘钥；jwtTTL 令牌有效期
// Returns: UserService 实例
func NewUserService(repo repository.UserRepository, jwtSecret []byte, jwtTTL time.Duration) UserService {
	return &userService{repo: repo, jwtSecret: jwtSecret, jwtTTL: jwtTTL}
}

// Register 注册用户并保存安全哈希密码
func (s *userService) Register(username, email, password string) (*user.User, error) {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}
	u := &user.User{Username: username, Email: email, PasswordHash: hash, Role: string(user.RoleUser), Status: "active"}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

// Login 校验凭据并生成访问令牌
func (s *userService) Login(email, password string) (string, *user.User, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if !auth.CheckPassword(password, u.PasswordHash) {
		return "", nil, errors.New("invalid credentials")
	}
	token, err := auth.CreateToken(u.ID, u.Role, s.jwtSecret, s.jwtTTL)
	if err != nil {
		return "", nil, err
	}
	return token, u, nil
}

// GetByID 查询用户信息
func (s *userService) GetByID(id uint) (*user.User, error) {
	return s.repo.GetByID(id)
}
