package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
