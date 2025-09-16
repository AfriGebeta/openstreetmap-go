package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateRandomPassword() string {
	var seed = time.Now().UnixNano()
	var random = rand.New(rand.NewSource(seed))

	var number = random.Int63n(89999999) + 10000000
	return fmt.Sprintf("%d", number)
}
