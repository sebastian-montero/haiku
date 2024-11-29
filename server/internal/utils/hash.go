package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, salt []byte) (string, error) {
	saltedPassword := append(salt, []byte(password)...)
	hashedPassword, err := bcrypt.GenerateFromPassword(saltedPassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
