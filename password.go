package gisproxy

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

//HashPassword return hash for given password
func HashPassword(password string) string {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(passHash)

}

//CheckPassword compare password and hashed password
func CheckPassword(password string, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		GetLogger().Error(err)
		return false
	}
	return true
}
