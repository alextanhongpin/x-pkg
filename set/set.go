// Package set implements a unique set of items.
package set

import (
	"fmt"
	"sync"
)

func Example() {
	s := New("a", "b", "c")
	s.Add("d", "e")
	fmt.Println(s.Has("d"))
	s.Remove("a", "b")
	fmt.Println(s.Has("a"))
	fmt.Println(s.Size())
}

type (
	// Set represents the set operations.
	Set interface {
		Has(key interface{}) bool
		Add(keys ...interface{})
		Remove(keys ...interface{})
		Size() int
	}

	// set implements the Set interface.
	set struct {
		value map[interface{}]struct{}
	}
)

// New creates a new set from the given values.
func New(values ...interface{}) *set {
	s := set{make(map[interface{}]struct{})}
	for _, v := range values {
		s.Add(v)
	}
	return &s
}

// Add adds an existing item to the set.
func (s *set) Add(keys ...interface{}) {
	for _, key := range keys {
		s.value[key] = struct{}{}
	}
}

// Has checks if the given item exists in the set.
func (s *set) Has(key interface{}) bool {
	_, exist := s.value[key]
	return exist
}

// Remove removes an item from the set.
func (s *set) Remove(keys ...interface{}) {
	for _, key := range keys {
		delete(s.value, key)
	}
}

// Size returns the number of items in the given set.
func (s *set) Size() int {
	return len(s.value)
}

// ConcurrentSet is a thread-safe Set implementation.
type ConcurrentSet struct {
	sync.RWMutex
	value map[interface{}]struct{}
}

// NewConcurrent returns a thread-safe Set implementation.
func NewConcurrent(values ...interface{}) *ConcurrentSet {
	c := ConcurrentSet{value: make(map[interface{}]struct{})}
	for _, v := range values {
		c.Add(v)
	}
	return &c
}

// Add adds an item to the set.
func (c *ConcurrentSet) Add(keys ...interface{}) {
	c.Lock()
	for _, key := range keys {
		c.value[key] = struct{}{}
	}
	c.Unlock()
}

// Has checks if a key exists in the set.
func (c *ConcurrentSet) Has(key interface{}) bool {
	c.RLock()
	_, exist := c.value[key]
	c.RUnlock()
	return exist
}

// Remove removes an item from the set.
func (c *ConcurrentSet) Remove(keys ...interface{}) {
	c.Lock()
	for _, key := range keys {
		delete(c.value, key)
	}
	c.Unlock()
}

// Size returns the size of the set.
func (c *ConcurrentSet) Size() int {
	c.RLock()
	size := len(c.value)
	c.RUnlock()
	return size
}
