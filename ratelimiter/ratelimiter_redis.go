package ratelimiter

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const script = `
-- ARGV[1]: The current timestamp in nanoseconds.
-- KEYS[1]: The key to rate limit, e.g. clientIP + userID/sessionID.

-- Delete all keys that are older than n nanoseconds ago.
redis.call('ZREMRANGEBYSCORE', KEYS[1], 0, ARGV[1] - %d)

-- Find the number of remaining tokens left.
if tonumber(redis.call('ZCARD', KEYS[1])) < %d then
	redis.call('ZADD', KEYS[1], ARGV[1], ARGV[1])
	return 'ok'
else
	return redis.error_reply('limit exceeded')
end
`

type Redis struct {
	client    *redis.Client
	frequency int
	duration  time.Duration
	script    string
}

func NewRedis(client *redis.Client, duration time.Duration, frequency int) *Redis {
	return &Redis{
		client:    client,
		frequency: frequency,
		duration:  duration,
		script:    fmt.Sprintf(script, duration, frequency),
	}
}

func (r *Redis) Allow(key string) bool {
	var (
		keys = []string{key}
		args = []interface{}{time.Now().UnixNano()}
	)
	res, err := r.client.Eval(r.script, keys, args...).Result()
	return res == "ok" && err == nil
}
