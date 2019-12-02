package gostrings

import (
	"crypto/rand"
	"encoding/base64"
	"log"
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

func Example() {
	b, err := RandBytes(10)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(b, len(b))

	s, err := RandString(10)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(s, len(s))
}
