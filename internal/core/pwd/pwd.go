package pwd

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"pubot/internal/core/config"
)

// Hash 哈希密码，返回格式 "salt$HashPasswd"
func Hash(passwd string) string {
	hash := sha256.New()
	hash.Write([]byte(config.Get().PwdSalt + passwd))
	hashedPasswd := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return fmt.Sprintf("%s$%s", config.Get().PwdSalt, hashedPasswd)
}

// Verify 密码认证
func Verify(storagePwd, inputPwd string) (bool, error) {
	// 分割存储的密码为盐和哈希值
	parts := strings.Split(storagePwd, "$")
	if len(parts) != 2 {
		return false, errors.New("invalid stored password format")
	}
	salt := parts[0]
	hashedPassword := parts[1]
	// 重新生成输入密码的哈希值
	inputHash := sha256.New()
	inputHash.Write([]byte(salt + inputPwd))
	inputHashedPassword := base64.URLEncoding.EncodeToString(inputHash.Sum(nil))
	// 使用 `subtle.ConstantTimeCompare` 进行时间安全的比较
	if subtle.ConstantTimeCompare([]byte(hashedPassword), []byte(inputHashedPassword)) == 1 {
		return true, nil
	}
	return false, nil
}
