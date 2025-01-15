package tests

import (
	"chrispaul1/chirpy/internal/auth"
	"testing"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	// Create token
	token, err := auth.MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("Error creating token: %v", err)
	}

	// Validate token
	validatedUserID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Error validating token: %v", err)
	}
	if validatedUserID != userID {
		t.Fatalf("Expected user ID %v, got %v", userID, validatedUserID)
	}
}
