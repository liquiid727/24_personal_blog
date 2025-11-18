// 用户仓储：封装用户数据的持久化访问
package repository

import (
    "gorm.io/gorm"
    "go_blog/internal/domain/user"
)

// UserRepository 定义用户相关的数据访问接口
type UserRepository interface {
    Create(u *user.User) error
    GetByEmail(email string) (*user.User, error)
    GetByID(id uint) (*user.User, error)
}

type userRepository struct {
    db *gorm.DB
}

// NewUserRepository 创建用户仓储实现
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

// Create 新增用户记录
func (r *userRepository) Create(u *user.User) error {
    return r.db.Create(u).Error
}

// GetByEmail 根据邮箱查询用户
func (r *userRepository) GetByEmail(email string) (*user.User, error) {
    var u user.User
    if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
        return nil, err
    }
    return &u, nil
}

// GetByID 根据主键查询用户
func (r *userRepository) GetByID(id uint) (*user.User, error) {
    var u user.User
    if err := r.db.First(&u, id).Error; err != nil {
        return nil, err
    }
    return &u, nil
}