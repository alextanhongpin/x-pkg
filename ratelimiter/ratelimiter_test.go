package ratelimiter_test

import (
	"fmt"
	"time"

	"github.com/alextanhongpin/pkg/ratelimiter"
)

func Example_in_memory() {
	var (
		frequency   = Per(time.Second, 5)
		burst       = 5
		inactiveTTL = time.Second
	)
	rl, cancel := ratelimiter.New(frequency, burst, inactiveTTL)
	defer cancel()

	// Use all quota.
	ok := rl.Allow("user_1")
	ok = rl.Allow("user_1")
	ok = rl.Allow("user_1")
	ok = rl.Allow("user_1")
	ok = rl.Allow("user_1")
	ok = rl.Allow("user_1")
	fmt.Println("allowed?", ok)

	// Sleep to recover.
	time.Sleep(2 * time.Second)
	ok = rl.Allow("user_1")
	fmt.Println("allowed?", ok)
}
