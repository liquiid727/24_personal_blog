// File: internal/auth/password.go
// Purpose: Password hashing and verification utilities using bcrypt.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: Provides secure hashing and comparison helpers for user credentials.
package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 生成密码哈希
// Input: 明文密码
// Output: bcrypt 哈希字符串或错误
func HashPassword(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CheckPassword 校验密码
// Params: p 明文；h 哈希
// Returns: 是否匹配
func CheckPassword(p, h string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p)) == nil
}
