package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Hash 生成密码哈希，自动生成盐
func Hash(password string) (string, error) {
	// bcrypt 会自动生成随机盐
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Verify 验证密码
func Verify(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
