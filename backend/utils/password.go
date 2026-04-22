package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword buat ngerubah "password123" jadi "$2a$12$..." (biar hacker pusing)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPassword buat nyocokin password pas login
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}