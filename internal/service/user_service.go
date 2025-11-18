// 用户服务：实现注册、登录与查询等领域逻辑
package service

import (
    "errors"
    "time"
    
    "go_blog/internal/auth"
    "go_blog/internal/domain/user"
    "go_blog/internal/repository"
)

// UserService 定义用户业务接口
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