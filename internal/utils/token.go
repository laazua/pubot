package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims 自定义 JWT Claims
type UserClaims struct {
	UserID   uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"` // "admin" 或 "user"
	jwt.RegisteredClaims
}

// JWTAuth 封装 JWT 相关操作
type JWTAuth struct {
	SecretKey []byte
	Expire    time.Duration
}

// NewJWTAuth 创建 JWTAuth 实例
func NewJWTAuth(secret string, expire time.Duration) *JWTAuth {
	return &JWTAuth{
		SecretKey: []byte(secret),
		Expire:    expire,
	}
}

// GenerateToken 生成 token
func (j *JWTAuth) GenerateToken(userID uint, username, role string) (string, error) {
	claims := &UserClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.Expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SecretKey)
}

// ParseToken 解析 token 并返回 UserClaims
func (j *JWTAuth) ParseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateToken 校验 token 是否有效
func (j *JWTAuth) ValidateToken(tokenString string) bool {
	_, err := j.ParseToken(tokenString)
	return err == nil
}

// GetUserFromToken 从 token 中获取用户信息
func (j *JWTAuth) GetUserFromToken(tokenString string) (*UserClaims, error) {
	return j.ParseToken(tokenString)
}
