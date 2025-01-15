package tests

import (
	"chrispaul1/chirpy/internal/auth"
	"net/http"
	"testing"
)

func TestGetBearerToke(t *testing.T) {

	header := make(http.Header)
	header.Add("Authorization", "Bearer yeet-m8")

	token, err := auth.GetBearerToken(header)
	if err != nil {
		t.Fatalf("Error, could not retireve token : %v", err)
	}

	if token != "yeet-m8" {
		t.Fatalf("Error, token is wrong :%s", token)
	}

}
