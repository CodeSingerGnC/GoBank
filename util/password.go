package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 返回 bcrpt 算法加密的密码 hash
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}	
	return string(hashedPassword), nil
}

// CheckPasswordHash 检查密码是否匹配
func CheckPasswordHash(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
