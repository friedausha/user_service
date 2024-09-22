package hash

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
	"sync"
)

var (
	cache = make(map[string]string)
	mu    sync.RWMutex
)

// EncryptPassword hashes a plain text password using bcrypt
func EncryptPassword(password string) (string, error) {
	// TODO: alternative default cost? why use it?
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with plain text password
func CheckPasswordHash(hash, password string) bool {
	mu.RLock()
	cachedPassword, found := cache[hash]
	//print(cachedPassword)
	//print(password)
	mu.RUnlock()
	if found && strings.TrimSpace(cachedPassword) == strings.TrimSpace(password) {
		//print("hash found")
		return true
	}
	err := bcrypt.CompareHashAndPassword([]byte(strings.TrimSpace(hash)), []byte(strings.TrimSpace(password)))
	if err == nil {
		// If comparison is successful, store in cache
		mu.Lock()
		cache[hash] = password
		mu.Unlock()
	}
	return err == nil
}

//
//// EncryptPassword hashes a plain text password using MD5
//func EncryptPassword(password string) (string, error) {
//	hash := md5.New()
//	_, err := hash.Write([]byte(password))
//	if err != nil {
//		return "", err
//	}
//
//	return hex.EncodeToString(hash.Sum(nil)), nil
//}

//
//// CheckPasswordHash compares a hashed password with plain text password
//func CheckPasswordHash(hash, password string) bool {
//	hashedPassword, err := EncryptPassword(password)
//	if err != nil {
//		return false
//	}
//	return strings.TrimSpace(hash) == strings.TrimSpace(hashedPassword)
//}
