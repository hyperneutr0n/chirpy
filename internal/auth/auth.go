package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPasword(password string) (string, error) {
	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("failed to create hash: %w", err)
	}

	return hashed, nil
}

func CheckHash(password, hashed string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hashed)
	if err != nil {
		return false, fmt.Errorf("failed comparing hashed: %w", err)
	}

	return match, nil	
}