// JWT 工具：创建与解析访问令牌
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

// CreateToken 创建 HS256 签名的 JWT 令牌，携带用户与角色信息
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

// ParseToken 解析令牌并返回 Claims
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