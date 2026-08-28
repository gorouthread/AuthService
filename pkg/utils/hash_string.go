package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashString(str string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash string: %w", err)
	}
	return string(hashedBytes), nil
}

func CheckString(str, hashedStr string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(str))
	if err != nil {
		return fmt.Errorf("invalid string: %w", err)
	}
	return nil
}
