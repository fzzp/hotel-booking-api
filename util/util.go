package util

import (
	"crypto/md5"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// MD5
func MD5(plainText string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(plainText)))
}

// Hash 加密密码
func Hash(plainText string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Matches 匹配密码是否正确
func Matches(plaintext string, hashPwsd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPwsd), []byte(plaintext))
}
