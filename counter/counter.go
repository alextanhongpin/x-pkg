// Package counter allows events to be tracked and expired, e.g. keeping track
// of the user's login error rate and blocking them from doing further task
// until the cooldown period is completed.

package counter

import (
	"sync"
	"time"
)

type Item struct {
	// The identifier of the policy.
	// id string
	value      int
	lastAccess time.Time
}

type Counter struct {
	sync.RWMutex
	data map[interface{}]*Item
	max  int
	// TODO: Now all the items will share the same ttl and max count. It is
	// preferable to create a separate policy for different events.
	ttl time.Duration
}

func New(max int, ttl time.Duration) *Counter {
	return &Counter{
		data: make(map[interface{}]*Item),
		max:  max,
		ttl:  ttl,
	}
}

// Put(Policy{min, max, path, identifier: clientIP})
// Should change this to put with key and value.
func (c *Counter) Increment(key interface{}) {
	c.Lock()
	it, ok := c.data[key]
	if !ok {
		it = &Item{value: 0}
		c.data[key] = it
	}
	it.value++
	it.lastAccess = time.Now()
	c.Unlock()
}

// type Policy interface {Allow(*Item)}
// type Policy struct {
//         ttl time.Duration
//         max int
// }
//
// func (p *Policy) Allow(value int, lastAccess time.Time) bool {
//         if value < p.max {
//                 return true
//         }
//         if time.Since(lastAccess) > p.ttl {
//                 return true
//         }
//         return false
// }

func (c *Counter) Allow(key interface{}) bool {
	c.Lock()
	defer c.Unlock()
	it, ok := c.data[key]
	if !ok {
		return true
	}
	if it.value < c.max {
		return true
	}
	if time.Since(it.lastAccess) > c.ttl {
		delete(c.data, key)
		return true
	}
	return false
}

// TODO: Add a cleanup method to ensure the expired keys are deleted.

// type Event struct {
//         Name     string
//         ClientIP string
// }

// func main() {
//         counter := New(3, 1*time.Second)
//         counter.Increment(Event{"john", "0.0.0.0"})
//         fmt.Println("is unblocked", counter.Allow(Event{"john", "0.0.0.0"}))
//
//         counter.Increment(Event{"john", "0.0.0.0"})
//         counter.Increment(Event{"john", "0.0.0.0"})
//         counter.Increment(Event{"john", "0.0.0.0"})
//
//         fmt.Println("is unblocked", counter.Allow(Event{"john", "0.0.0.0"}))
//         time.Sleep(2 * time.Second)
//         fmt.Println("is unblocked", counter.Allow(Event{"john", "0.0.0.0"}))
// }
