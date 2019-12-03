package ttlmap

import (
	"sync"
	"time"
)

type Map interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	SetEx(key string, value interface{}, duration time.Duration)
}

type item struct {
	value     interface{}
	updatedAt time.Time
	ttl       time.Duration
}

type TTLMap struct {
	sync.RWMutex
	values map[string]item

	sync.Once
	quit chan interface{}

	wg sync.WaitGroup
}

func New() (*TTLMap, func()) {
	ttl := &TTLMap{
		values: make(map[string]item),
		quit:   make(chan interface{}),
	}
	go ttl.clear(5 * time.Second)
	return ttl, ttl.cancel
}

func (t *TTLMap) Set(key string, value interface{}) {
	t.Lock()
	// Negative ttl means forever.
	t.values[key] = item{value: value, updatedAt: time.Now(), ttl: -1}
	t.Unlock()
}

func (t *TTLMap) Get(key string) (interface{}, bool) {
	t.RLock()
	res, ok := t.values[key]
	t.RUnlock()

	if !ok {
		return nil, ok
	}
	t.Lock()
	if res.ttl > 0 && time.Since(res.updatedAt) > res.ttl {
		delete(t.values, key)
		return nil, false
	}
	t.Unlock()
	return res.value, true
}

func (t *TTLMap) SetEx(key string, value interface{}, duration time.Duration) {
	t.Lock()
	t.values[key] = item{value: value, updatedAt: time.Now(), ttl: duration}
	t.Unlock()
}

func (t *TTLMap) clear(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	t.wg.Add(1)
	defer t.wg.Done()
	for {
		select {
		case <-t.quit:
			return
		case <-ticker.C:
			t.Lock()
			var i int
			for key, item := range t.values {
				if item.ttl > 0 && time.Since(item.updatedAt) > item.ttl {
					i++
					delete(t.values, key)
				}
				// Clean at most 10,000 keys to prevent holding
				// the lock for too long.
				if i > 10_000 {
					break
				}
			}
			t.Unlock()
		}
	}
}

func (t *TTLMap) cancel() {
	t.Once.Do(func() {
		close(t.quit)
		t.wg.Wait()
	})
}
