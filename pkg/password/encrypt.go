package password

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func Hash(password []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		return ""
	}
	return string(hashed)
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {

	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
