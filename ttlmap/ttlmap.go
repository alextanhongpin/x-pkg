package ttlmap

import (
	"context"
	"sync"
	"time"
)

type item struct {
	value      string
	lastAccess time.Time
	// lastAccess int64
}

type TTLMap struct {
	sync.RWMutex
	data map[string]*item
	// Good for performance, but not robust against wall clock reset.
	// ttl  int64
	ttl  time.Duration
	once sync.Once
	wg   sync.WaitGroup
	quit chan interface{}
}

func New(ttl time.Duration) *TTLMap {
	return &TTLMap{
		data: make(map[string]*item),
		quit: make(chan interface{}),
		ttl:  ttl,
	}
}

func (t *TTLMap) Cleanup(every time.Duration) func(context.Context) {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		t.cleanup(every)
	}()
	return func(ctx context.Context) {
		done := make(chan interface{})
		t.once.Do(func() {
			close(t.quit)
			t.wg.Wait()
			close(done)
		})

		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		}
	}
}

func (t *TTLMap) Len() int {
	t.RLock()
	l := len(t.data)
	t.RUnlock()
	return l
}

func (t *TTLMap) Put(key, value string) {
	t.Lock()
	it, ok := t.data[key]
	if !ok {
		it = &item{value: value}
		t.data[key] = it
	}
	it.lastAccess = time.Now()
	t.Unlock()
}

func (t *TTLMap) Get(key string) (value string) {
	t.Lock()
	if it, ok := t.data[key]; ok {
		value = it.value
		it.lastAccess = time.Now()
	}
	t.Unlock()
	return
}

func (t *TTLMap) cleanup(every time.Duration) {
	ticker := time.NewTicker(every)
	defer ticker.Stop()
	for {
		select {
		case <-t.quit:
			return
		case <-ticker.C:
			t.Lock()
			for k, v := range t.data {
				if time.Since(v.lastAccess) > t.ttl {
					// if time.Now().Unix()-v.lastAccess > t.ttl {
					// log.Println("deleted", k)
					delete(t.data, k)
				}
			}
			t.Unlock()
		}
	}
}
