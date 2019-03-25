package observable

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type AsyncObservable struct {
	wg sync.WaitGroup
	// We can probably extract this into a separate data structure.
	sync.RWMutex
	events map[Event][]Action

	// Ensure the write channel is closed once.
	once sync.Once
	quit chan interface{}

	// Ensure the read channel is closed once.
	onceCh sync.Once
	ch     chan Message
}

func NewAsync(n int) *AsyncObservable {
	events := make(map[Event][]Action)
	return &AsyncObservable{
		events: events,
		quit:   make(chan interface{}),
		ch:     make(chan Message, n),
	}
}

func (o *AsyncObservable) On(event Event, fn Action) {
	o.Lock()
	_, exist := o.events[event]
	if !exist {
		o.events[event] = make([]Action, 0)
	}
	o.events[event] = append(o.events[event], fn)
	o.Unlock()
}

func (o *AsyncObservable) Emit(event Event, params interface{}) error {
	// Why do we need two select here? This is to ensure only the sole case here (o.quit) does not compete with the o.ch. In a situation where both o.quit and o.ch matches the select, if the o.ch is closed and a message is sent to o.ch, it will panic.
	select {
	case <-o.quit:
		o.onceCh.Do(func() {
			close(o.ch)
		})
		return errors.New("channel closed")
	default:
	}

	select {
	case <-o.quit:
		return errors.New("channel closed")
	case o.ch <- Message{event, params}:
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("timeout exceeded")
	}
}

// Stop terminates all the write channels first before closing the read
// channels. The unread messages will be flushed before the process completes.
func (o *AsyncObservable) Stop() {
	o.once.Do(func() {
		close(o.quit)
	})
	o.wg.Wait()
}

// Start
func (o *AsyncObservable) Start() {
	o.wg.Add(1)
	go func() {
		defer o.wg.Done()
		for evt := range o.ch {
			// Still need to lock to ensure there's no data race.
			o.RLock()
			fns, exist := o.events[evt.event]
			o.RUnlock()
			if !exist {
				log.Println(fmt.Errorf(`event "%s" is not registered`, evt.event))
			}
			for _, fn := range fns {
				if err := fn(evt.params); err != nil {
					log.Println(err)
				}
			}
		}
	}()

}

const (
	LoginEvent   = Event("login")
	LogoutEvent  = Event("logout")
	UnknownEvent = Event("unknown")
)

type UserService struct {
	Observer
}

func (u *UserService) Unknown() {
	_ = u.Emit(UnknownEvent, nil)
}
func (u *UserService) Logout(msg string) {
	_ = u.Emit(LogoutEvent, msg)
}

type LoginRequest struct {
	Msg string
}

func (u *UserService) Login(req LoginRequest) {
	sent := u.Emit(LoginEvent, req)
	_ = sent
}
