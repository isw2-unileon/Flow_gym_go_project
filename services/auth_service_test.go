package services

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "my_secure_password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if hash == "" {
		t.Error("Expected hashed password, got empty string")
	}
	if hash == password {
		t.Error("Expected hash to be different from plain password")
	}

	tooLongPassword := strings.Repeat("a", 73)
	_, err = HashPassword(tooLongPassword)
	if err == nil {
		t.Error("Expected an error for password > 72 bytes, got nil")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "my_secure_password"
	hash, _ := HashPassword(password)

	if !CheckPassword(password, hash) {
		t.Error("Expected password to match hash")
	}

	if CheckPassword("wrong_password", hash) {
		t.Error("Expected password to NOT match hash")
	}
}