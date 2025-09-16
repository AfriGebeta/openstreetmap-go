package v1

import (
	"encoding/hex"
	"math/rand"
)

func GenerateRandomHex(byteLength int) (string, error) {
	bytes := make([]byte, byteLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
