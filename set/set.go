package set

import (
	"sync"
)

type (
	Set interface {
		Has(key interface{}) bool
		Add(key interface{})
		Remove(key interface{})
		Size() int
	}
	SetImpl struct {
		value map[interface{}]struct{}
	}
)

func New() *SetImpl {
	return &SetImpl{make(map[interface{}]struct{})}
}

func (s *SetImpl) Add(key interface{}) {
	s.value[key] = struct{}{}
}

func (s *SetImpl) Has(key interface{}) bool {
	_, exist := s.value[key]
	return exist
}

func (s *SetImpl) Remove(key interface{}) {
	delete(s.value, key)
}

func (s *SetImpl) Size() int {
	return len(s.value)
}

type ConcurrentSet struct {
	sync.RWMutex
	value map[interface{}]struct{}
}

func NewConcurrent() *ConcurrentSet {
	value := make(map[interface{}]struct{})
	return &ConcurrentSet{value: value}
}

func (c *ConcurrentSet) Add(key interface{}) {
	c.Lock()
	c.value[key] = struct{}{}
	c.Unlock()
}

func (c *ConcurrentSet) Has(key interface{}) bool {
	c.RLock()
	_, exist := c.value[key]
	c.RUnlock()
	return exist
}

func (c *ConcurrentSet) Remove(key interface{}) {
	c.Lock()
	delete(c.value, key)
	c.Unlock()
}

func (c *ConcurrentSet) Size() int {
	c.RLock()
	size := len(c.value)
	c.RUnlock()
	return size
}
