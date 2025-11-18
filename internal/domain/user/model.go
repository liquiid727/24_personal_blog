// 用户领域模型：管理注册用户与角色
package user

import "time"

// Role 用户角色类型
type Role string

// 预定义角色：管理员与普通用户
const (
    RoleAdmin Role = "admin"
    RoleUser  Role = "user"
)

// User 表示系统中的注册用户
type User struct {
    ID           uint      `gorm:"primaryKey"`
    Username     string    `gorm:"uniqueIndex;size:64"`   // 唯一用户名
    Email        string    `gorm:"uniqueIndex;size:128"`  // 唯一邮箱
    PasswordHash string    `gorm:"size:256"`              // bcrypt 哈希
    Role         string    `gorm:"size:16"`               // 角色（admin/user）
    Status       string    `gorm:"size:16"`               // 状态（active/disabled）
    CreatedAt    time.Time
    UpdatedAt    time.Time
}