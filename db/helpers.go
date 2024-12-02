package db

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func IsPasswordHashed(password string) bool {
	return strings.HasPrefix(password, "$2a$")
}
