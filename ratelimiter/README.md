# RateLimiter


```go
lim := ratelimiter.New(ratelimiter.Per(1 * time.Second, 3))
shutdown := lim.CleanupVisitor(5 * time.Second, 3 * time.Second)
```
