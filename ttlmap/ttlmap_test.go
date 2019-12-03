package ttlmap_test

import (
	"log"
	"time"

	"github.com/alextanhongpin/pkg/ttlmap"
)

func Example() {
	m, cancel := ttlmap.New()
	defer cancel()

	var (
		key      = "hello"
		value    = 1
		duration = time.Second
	)
	m.Set(key, value)
	m.Get(key)
	m.SetEx(key, value, duration)

	val, ok := m.Get(key)
	if !ok {
		log.Fatal("key not present")
	}
	log.Println("got val", val)

	time.Sleep(2 * time.Second)
	val, ok = m.Get(key)
	if !ok {
		log.Fatal("key not present")
	}
	log.Println("got val", val)
}
