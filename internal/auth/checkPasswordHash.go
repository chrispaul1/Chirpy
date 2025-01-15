package auth

import "golang.org/x/crypto/bcrypt"

func CheckPasswordHash(password, hash string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
