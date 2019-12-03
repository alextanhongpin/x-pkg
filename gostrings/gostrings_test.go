package gostrings_test

import (
	"log"

	"github.com/alextanhongpin/pkg/gostrings"
)

func Example() {
	b, err := gostrings.RandBytes(10)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(b, len(b))

	s, err := gostrings.RandString(10)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(s, len(s))
}
