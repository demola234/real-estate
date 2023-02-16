package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// Hashing the users password with bcrypt
	bytes, error := bcrypt.GenerateFromPassword([]byte(password), 14)
	if error != nil {
		log.Panic(error)
	}

	return string(bytes), error
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "Invalid password"
		check = false
	}
	return check, msg
}
