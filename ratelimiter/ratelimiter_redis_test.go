package ratelimiter_test

import (
	"fmt"
	"time"

	"github.com/alextanhongpin/pkg/ratelimiter"

	"github.com/go-redis/redis/v8"
)

func Example_redis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	var (
		userID = "127.0.0.1/123"
	)
	// 5 requests per second.
	rl := ratelimiter.NewRedis(client, time.Second, 5)
	for i := 0; i < 10; i++ {
		ok := rl.Allow(userID)
		fmt.Println(ok)
	}
}
