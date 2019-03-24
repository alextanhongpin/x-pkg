## Synchronous Observable

```go
package main

import (
	"log"
	"time"

	"github.com/alextanhongpin/pkg/observable"
)

type Service struct {
	observable.Observer
}

func (s *Service) Greet(msg string) {
	s.Emit("greeted", msg)

	log.Println("service.Greet called")
}

func main() {
	svc := &Service{observable.NewSyncObservable()}
	svc.On("greeted", func(msg interface{}) error {
		log.Println("got:", msg)
		time.Sleep(1 * time.Second)
		return nil
	})
	svc.Greet("hello")
}
```

## Asynchronous Observable

```go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Observable struct {
	event ObserverEvent
	fn    ObserverFunc
}

type ObserverEvent string
type ObserverFunc func(interface{}) error
type Observer interface {
	On(event ObserverEvent, fn ObserverFunc) bool
	Emit(event ObserverEvent, params interface{}) bool
	Start()
	Stop()
}

type ObserverMessage struct {
	event  ObserverEvent
	params interface{}
}

type ObserverImpl struct {
	wg sync.WaitGroup

	events   map[ObserverEvent][]ObserverFunc
	register chan Observable

	once sync.Once
	quit chan interface{}
	ch   chan ObserverMessage
}

func NewObserver(n int) *ObserverImpl {
	events := make(map[ObserverEvent][]ObserverFunc)
	return &ObserverImpl{
		events:   events,
		register: make(chan Observable),
		quit:     make(chan interface{}),
		ch:       make(chan ObserverMessage, n),
	}
}

func (o *ObserverImpl) On(event ObserverEvent, fn ObserverFunc) bool {
	select {
	case <-o.quit:
		return false
	case o.register <- Observable{event, fn}:
		return true
	case <-time.After(5 * time.Second):
		return false
	}

}

func (o *ObserverImpl) Emit(event ObserverEvent, params interface{}) bool {
	select {
	case <-o.quit:
		return false
	case o.ch <- ObserverMessage{event, params}:
		return true
	case <-time.After(5 * time.Second):
		return false
	}
}

func (o *ObserverImpl) Stop() {
	o.once.Do(func() {
		close(o.quit)
	})
	o.wg.Wait()
	log.Println("stopped")
}

func (o *ObserverImpl) Start() {
	o.wg.Add(1)
	go func() {
		defer o.wg.Done()
		for {
			select {
			case <-o.quit:
				return
			case obs, ok := <-o.register:
				if !ok {
					return
				}
				_, exist := o.events[obs.event]
				if !exist {
					o.events[obs.event] = make([]ObserverFunc, 0)
				}
				o.events[obs.event] = append(o.events[obs.event], obs.fn)
			case evt, ok := <-o.ch:
				if !ok {
					return
				}
				fns, exist := o.events[evt.event]
				if !exist {
					log.Println(fmt.Errorf(`event "%s" is not registered`, evt.event))
				}
				for _, fn := range fns {
					// time.Sleep(1 * time.Second)
					if err := fn(evt.params); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}()

}

const (
	LoginEvent   = ObserverEvent("login")
	LogoutEvent  = ObserverEvent("logout")
	UnknownEvent = ObserverEvent("unknown")
)

type UserService struct {
	Observer
}

func (u *UserService) Unknown() {
	_ = u.Emit(UnknownEvent, nil)
}
func (u *UserService) Logout(msg string) {
	fmt.Println("logging out", msg)
	sent := u.Emit(LogoutEvent, msg)
	log.Println("sent logout", sent)
}

type LoginRequest struct {
	Msg string
}

func (u *UserService) Login(req LoginRequest) {
	// Do some work...
	sent := u.Emit(LoginEvent, req)
	log.Println("sent login", sent)
}

func main() {
	usersvc := &UserService{NewObserver(10)}
	usersvc.Start()

	usersvc.On(LogoutEvent, func(msg interface{}) error {
		fmt.Println("on logout:", msg)
		return nil
	})
	usersvc.Logout("john has logged out")

	usersvc.On(LoginEvent, func(msg interface{}) error {
		switch v := msg.(type) {
		case LoginRequest:
			fmt.Println("on login:", v.Msg)
		default:
			fmt.Println("not handled")
		}
		return nil
	})

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			usersvc.Login(LoginRequest{fmt.Sprintf("user %d", i+1)})
			usersvc.Unknown()
		}
	}()

	time.Sleep(4 * time.Second)
	usersvc.Stop()

	time.Sleep(3 * time.Second)
}
```
