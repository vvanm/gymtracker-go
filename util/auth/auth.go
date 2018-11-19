package auth

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func SecurePassword(p string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
func PasswordMatch(in string, out string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(out), []byte(in))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
