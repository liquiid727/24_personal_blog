// File: internal/auth/jwt.go
// Purpose: Utilities to create and parse JWT access tokens.
// Author: Go Blog Team
// Created: 2025-11-18
// Last Modified: 2025-11-18
// Description: HS256-signed tokens carrying user ID and role; helpers for issuance and validation.
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义声明，包含用户 ID 与角色
type Claims struct {
	UID  uint   `json:"uid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// CreateToken 创建令牌
// Params: uid 用户ID；role 角色；secret HS256秘钥；ttl 过期时长
// Returns: 签名后的 token 字符串或错误
func CreateToken(uid uint, role string, secret []byte, ttl time.Duration) (string, error) {
	now := time.Now()
	c := Claims{
		UID:  uid,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(secret)
}

// ParseToken 解析令牌
// Params: tokenStr 令牌文本；secret HS256秘钥
// Returns: *Claims 若校验通过，否则错误
// Errors: 包含签名错误、过期、格式不合法等
func ParseToken(tokenStr string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := token.Claims.(*Claims); ok && token.Valid {
		return c, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
