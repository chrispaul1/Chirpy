package tests

import (
	"chrispaul1/chirpy/internal/auth"
	"testing"
)

// Testing is hashPassword works
func TestHashPas(t *testing.T) {
	passwordToHash := "StableHorseBattery"
	_, err := auth.HashPassword(passwordToHash)
	if err != nil {
		t.Fatalf("Error, password '%s' did not hash", passwordToHash)
	}
}

// Test's if the check password hash func can detect the same hashed password
func TestCheckPasswordHashPass(t *testing.T) {

	password1 := "helloworld"
	password2 := "HelloSworld"

	hashedPassword1, err := auth.HashPassword(password1)
	if err != nil {
		t.Fatalf("Error, password '%s' did not hash", hashedPassword1)
	}

	hashedPassword2, err := auth.HashPassword(password2)
	if err != nil {
		t.Fatalf("Error, password '%s' did not hash", hashedPassword2)
	}

	stringHash1 := string(hashedPassword1)
	stringHash2 := string(hashedPassword2)
	err = auth.CheckPasswordHash(password1, stringHash1)
	if err != nil {
		t.Fatalf("Error, password should have matched\npassword : %s\nhash\t : %s", stringHash1, stringHash1)
	}

	err = auth.CheckPasswordHash(password2, stringHash1)
	if err == nil {
		t.Fatalf("Error, password should not have matche\n password: %s\nhash : %s", stringHash1, stringHash2)
	}
}
