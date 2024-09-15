package hash

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// EncryptPassword hashes a plain text password using bcrypt
func EncryptPassword(password string) (string, error) {
	// TODO: alternative default cost? why use it?
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with plain text password
func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(strings.TrimSpace(hash)), []byte(strings.TrimSpace(password)))
	return err == nil
}
