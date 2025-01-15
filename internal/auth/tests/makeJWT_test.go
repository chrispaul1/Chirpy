package tests

import (
	"chrispaul1/chirpy/internal/auth"
	"testing"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	testID := uuid.New()
	tokenSecret := "RocketRacoon"
	token, err := auth.MakeJWT(testID, tokenSecret)
	if err != nil {
		t.Fatalf("Error in creating the token : %v", err)
	}

	if token == "" {
		t.Fatalf(("Error, token is empty"))
	}

}
