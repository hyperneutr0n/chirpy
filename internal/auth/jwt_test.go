package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	duration := time.Hour

	token, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	if token == "" {
		t.Fatal("Expected non-empty token")
	}
}

func TestValidateJWT_Success(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	duration := time.Hour

	token, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	validatedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if validatedID != userID {
		t.Errorf("Expected user ID %v, got %v", userID, validatedID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	duration := -time.Hour

	token, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("Expected error for expired token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	wrongSecret := "wrong-secret"
	duration := time.Hour

	token, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("Expected error for wrong secret, got nil")
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	secret := "test-secret"
	invalidToken := "not.a.valid.jwt"

	_, err := ValidateJWT(invalidToken, secret)
	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}
}
