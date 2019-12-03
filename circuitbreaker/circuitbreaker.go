// Package circuitbreaker implements an in-memory circuit breaker to add
// resiliency to single instance.
package circuitbreaker

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ErrTooManyRequests is returned when the CircuitBreaker is in open state.
var ErrTooManyRequests = errors.New("too many requests")

// Task represents the task the circuitbreaker will be executing, and returns
// either a response or error.
type Task func() (interface{}, error)

// State mantains the circuit breaker state and is concurrent-safe.
type State struct {
	sync.RWMutex
	FailureCounter   int
	FailureThreshold int
	SuccessCounter   int
	SuccessThreshold int
	Timeout          time.Duration
	Timer            time.Time
}

// StartTimeout starts the timeout timer.
func (s *State) StartTimeoutTimer() {
	s.Lock()
	s.Timer = time.Now()
	s.Unlock()
}

// IncrementSuccessCounter increments the success counter by 1.
func (s *State) IncrementSuccessCounter() {
	s.Lock()
	s.SuccessCounter++
	s.Unlock()
}

// ResetSuccessCounter reset the success counter back to 0.
func (s *State) ResetSuccessCounter() {
	s.Lock()
	s.SuccessCounter = 0
	s.Unlock()
}

// IncrementFailureCounter increments the failure counter by 1.
func (s *State) IncrementFailureCounter() {
	s.Lock()
	s.FailureCounter++
	s.Unlock()
}

// ResetFailureCounter reset the failure counter back to 0.
func (s *State) ResetFailureCounter() {
	s.Lock()
	s.FailureCounter = 0
	s.Unlock()
}

// IsTimeoutTimerExpired checks if the timeout timer has expired.
func (s *State) IsTimeoutTimerExpired() bool {
	s.RLock()
	timer, timeout := s.Timer, s.Timeout
	s.RUnlock()
	return time.Since(timer) > timeout
}

// IsFailureThresholdExceeded checks if the failure threshold has exceed.
func (s *State) IsFailureThresholdExceeded() bool {
	s.RLock()
	failureCounter, failureThreshold := s.FailureCounter, s.FailureThreshold
	s.RUnlock()
	return failureCounter > failureThreshold
}

// IsSuccessThresholdExceeded checks if the success threshold has exceed.
func (s *State) IsSuccessThresholdExceeded() bool {
	s.RLock()
	successCounter, successThreshold := s.SuccessCounter, s.SuccessThreshold
	s.RUnlock()
	return successCounter > successThreshold
}

// CircuitBreaker represents the state machine for the circuit breaker
// algorithm.
type CircuitBreaker interface {
	Next() CircuitBreaker
	Handle(Task) (interface{}, error)
}

// Closed represents the closed state.
type Closed struct {
	state *State
}

// NewClosed returns a new Closed state.
func NewClosed(state *State) *Closed {
	// entry/reset failure counter
	state.ResetFailureCounter()
	return &Closed{state}
}

// Next checks if the next state transition is possible, and returns the next
// state or itself.
func (c *Closed) Next() CircuitBreaker {
	// failure threshold reached
	if c.state.IsFailureThresholdExceeded() {
		fmt.Println("is opened")
		return NewOpened(c.state)
	}
	return c
}

// Handle takes a Task and executes it on its behalf.
func (c *Closed) Handle(task Task) (interface{}, error) {
	// do/	if operation succeeds
	// 		return result
	// 	else
	// 		increment failure counter
	//		return failure
	res, err := task()
	if err != nil {
		c.state.IncrementFailureCounter()
		return nil, err
	}
	return res, nil
}

// Opened represents a new opened state.
type Opened struct {
	state *State
}

// NewOpened returns a new Opened state.
func NewOpened(state *State) *Opened {
	// entry/ start timeout timer
	state.StartTimeoutTimer()
	return &Opened{state}
}

// Next checks if the next state transition is possible, and returns the next
// state or itself.
func (o *Opened) Next() CircuitBreaker {
	// timeout timer expired
	if o.state.IsTimeoutTimerExpired() {
		fmt.Println("is half-opened")
		return NewHalfOpened(o.state)
	}
	return o
}

// Handle takes a Task and executes it on its behalf.
func (o *Opened) Handle(task Task) (interface{}, error) {
	// do /return failure
	return nil, ErrTooManyRequests
}

// HalfOpened represents the half opened state.
type HalfOpened struct {
	state  *State
	failed int32
}

// NewHalfOpened returns a new half-opened state.
func NewHalfOpened(state *State) *HalfOpened {
	// entry/ reset success counter
	state.ResetSuccessCounter()
	return &HalfOpened{state: state}
}

// Handle takes a Task and executes it on its behalf.
func (h *HalfOpened) Handle(task Task) (interface{}, error) {
	// do/	if operation succeeds
	// 		increment success counter
	// 		return result
	// 	else
	// 		return failure
	res, err := task()
	if err != nil {
		atomic.CompareAndSwapInt32(&h.failed, 0, 1)
		return nil, err
	}
	atomic.CompareAndSwapInt32(&h.failed, 1, 0)
	h.state.IncrementSuccessCounter()
	return res, err
}

// Next checks if the next state transition is possible, and returns the next
// state or itself.
func (h *HalfOpened) Next() CircuitBreaker {
	// success count threshold reached
	if h.state.IsSuccessThresholdExceeded() {
		fmt.Println("success count threshold reached")
		fmt.Println("is closed")
		return NewClosed(h.state)
	}
	if atomic.LoadInt32(&h.failed) == 1 {
		// operation failed
		fmt.Println("operation failed")
		return NewOpened(h.state)
	}

	return h
}

// CircuitBreakerImpl implements the CircuitBreaker interface.
type CircuitBreakerImpl struct {
	CircuitBreaker
}

// NewDefaultState returns a default state for the circuit breaker.
func NewDefaultState() *State {
	return &State{
		FailureCounter:   5,
		FailureThreshold: 5,
		SuccessCounter:   5,
		SuccessThreshold: 5,
		Timeout:          5 * time.Second,
	}
}

// New returns a new pointer to the CircuitBreaker implementation.
func New(state *State) *CircuitBreakerImpl {
	if state == nil {
		state = NewDefaultState()
	}
	cb := NewClosed(state)
	return &CircuitBreakerImpl{cb}
}

// Handle decorates the state.
func (c *CircuitBreakerImpl) Handle(task Task) (interface{}, error) {
	c.CircuitBreaker = c.Next()
	return c.CircuitBreaker.Handle(task)
}
