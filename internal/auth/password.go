// 密码工具：提供安全哈希与校验
package auth

import (
    "golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 生成密码哈希
func HashPassword(p string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(b), nil
}

// CheckPassword 对比明文密码与哈希是否匹配
func CheckPassword(p, h string) bool {
    return bcrypt.CompareHashAndPassword([]byte(h), []byte(p)) == nil
}