package password

import (
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

func BcryptPassword(plainPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 11)
	return string(bytes), err
}

func BcryptCompare(passwordHash, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plainPassword))
	return err == nil
}

func BcryptGetCost(passwordHash string) int {
	cost, _ := bcrypt.Cost([]byte(passwordHash))
	return cost
}
