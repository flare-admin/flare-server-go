package random

import (
	"crypto/rand"
	"math/big"
)

const (
	passwordLength = 8
	characters     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateDefPwd() string {
	if pass, err := GeneratePassword(passwordLength); err != nil {
		return "12345678"
	} else {
		return pass
	}
}

func GeneratePassword(length int) (string, error) {
	password := make([]byte, length)
	for i := range password {
		charIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		if err != nil {
			return "", err
		}
		password[i] = characters[charIndex.Int64()]
	}
	return string(password), nil
}
