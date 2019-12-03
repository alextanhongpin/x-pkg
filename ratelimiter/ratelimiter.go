// Package ratelimiter implements an in-memory rate limiter.
package ratelimiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	updatedAt time.Time
	limiter   *rate.Limiter
}

type RateLimiter struct {
	sync.RWMutex
	clients map[string]client

	factory func() *rate.Limiter
	sync.Once
	quit chan interface{}
	wg   sync.WaitGroup
}

func Per(interval time.Duration, times int) rate.Limit {
	frequency := interval / time.Duration(times)
	return rate.Every(frequency)
}

func New(frequency rate.Limit, burst int, inactiveTTL time.Duration) (*RateLimiter, func()) {
	rateLimiter := &RateLimiter{
		quit:    make(chan interface{}),
		clients: make(map[string]client),
		factory: func() *rate.Limiter {
			return rate.NewLimiter(frequency, burst)
		},
	}
	go rateLimiter.clean(inactiveTTL)
	return rateLimiter, rateLimiter.cancel
}

func (r *RateLimiter) clean(inactiveTTL time.Duration) {
	ticker := time.NewTicker(inactiveTTL * 2)
	defer ticker.Stop()

	r.wg.Add(1)
	defer r.wg.Done()

	for {
		select {
		case <-r.quit:
			return
		case <-ticker.C:
			r.Lock()
			var i int
			for id, client := range r.clients {
				if time.Since(client.updatedAt) > inactiveTTL {
					i++
					delete(r.clients, id)
				}

				// Clean at most 10,000 keys to prevent holding
				// the lock for too long.
				if i > 10_000 {
					break
				}
			}
			r.Unlock()
		}
	}
}

func (r *RateLimiter) Allow(key string) bool {
	client := r.get(key)
	return client.limiter.Allow()
}

func (r *RateLimiter) get(key string) client {
	r.RLock()
	client, ok := r.clients[key]
	r.RUnlock()

	if !ok {
		return r.add(key)
	}
	r.update(key)
	return client
}

func (r *RateLimiter) add(key string) client {
	c := client{
		limiter:   r.factory(),
		updatedAt: time.Now(),
	}
	r.Lock()
	r.clients[key] = c
	r.Unlock()
	return c
}

func (r *RateLimiter) update(key string) {
	r.Lock()
	if client, ok := r.clients[key]; ok {
		client.updatedAt = time.Now()
	}
	r.Unlock()
}

func (r *RateLimiter) cancel() {
	r.Once.Do(func() {
		close(r.quit)
		r.wg.Wait()
	})
}
