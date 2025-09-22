package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

func GenerateConfirmationToken(userID int64, secretKey string) (string, error) {
	payload := fmt.Sprintf("%d:%d", userID, time.Now().Unix())

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(payload))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	token := fmt.Sprintf("%s--%s", signature, base64.URLEncoding.EncodeToString([]byte(payload)))

	return token, nil
}
