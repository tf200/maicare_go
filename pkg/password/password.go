package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(password, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashed), nil
}
