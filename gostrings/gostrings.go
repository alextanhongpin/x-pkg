package gostrings

import (
	"crypto/rand"
	"encoding/base64"
)

func RandBytes(n int) ([]byte, error) {
	if n <= 0 {
		n = 8
	}
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func RandString(n int) (string, error) {
	b, err := RandBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}
