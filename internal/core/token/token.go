package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"pubot/internal/core/config"
)

type payLoad struct {
	Name      string
	ExpiresAt int64
}

// Create 创建token
func Create(name string) (string, error) {
	payload := payLoad{
		Name:      name,
		ExpiresAt: time.Now().Add(config.Get().TokenExpired * time.Minute).Unix(),
	}
	payLoadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	// 对 payload 进行 Base64 编码
	encodedPayload := base64.URLEncoding.EncodeToString(payLoadBytes)
	// 生成签名
	signature := generateSignature(encodedPayload, config.Get().SecretKey)
	// 返回 token
	return encodedPayload + "." + signature, nil
}

// Parse 解析token
func Parse(token string) (*payLoad, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}
	encodedPayload := parts[0]
	signature := parts[1]
	// 验证签名
	expectedSignature := generateSignature(encodedPayload, config.Get().SecretKey)
	if signature != expectedSignature {
		return nil, errors.New("invalid token signature")
	}
	// 解码 payload
	payloadBytes, err := base64.URLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return nil, errors.New("failed to decode payload")
	}
	// 解析 JSON
	var payload payLoad
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return nil, errors.New("failed to parse payload")
	}
	// 检查是否过期
	if payload.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token has expired")
	}
	return &payload, nil
}

// generateSignature 生成签名
func generateSignature(payload string, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(payload))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
