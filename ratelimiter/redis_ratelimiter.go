package ratelimiter

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const script = `
-- ARGV[1]: The current timestamp in seconds.
-- KEYS[1]: The key to rate limit, e.g. clientIP + userID/sessionID.

-- Delete all keys that are older than n seconds ago.
redis.call('ZREMRANGEBYSCORE', KEYS[1], 0, ARGV[1] - %d)

-- Find the number of remaining tokens left.
if tonumber(redis.call('ZCARD', KEYS[1])) < %d then
	redis.call('ZADD', KEYS[1], ARGV[1], ARGV[1])
	return 'ok'
else
	return redis.error_reply('limit exceeded')
end
`

type RedisRateLimiter struct {
	client    *redis.Client
	perSecond int
}

func NewRedisRateLimiter(client *redis.Client, perSecond int) *RedisRateLimiter {
	return &RedisRateLimiter{client, perSecond}
}

func (r *RedisRateLimiter) Allow(key string) bool {
	res, err := r.client.Eval(fmt.Sprintf(script, time.Second, r.perSecond), []string{key}, time.Now().UnixNano()).Result()
	return res == "ok" && err == nil
}

func Example() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	var (
		userID    = "1"
		perSecond = 5
	)
	// 5 requests per second.
	rl := NewRedisRateLimiter(client, perSecond)
	for i := 0; i < 10; i++ {
		ok := rl.Allow(userID)
		fmt.Println(ok)
	}
}
