## Synchronous Observable

```go
package main

import (
	"log"
	"time"

	"github.com/alextanhongpin/pkg/observable"
)

var GreetEvent = observable.Event("greeted")

type Service struct {
	observable.Observer
}

func (s *Service) Greet(msg string) {
	s.Emit(GreetEvent, msg)

	log.Println("service.Greet called")
}

func main() {
	svc := &Service{observable.NewSync()}
	svc.On(GreetEvent, func(msg interface{}) error {
		time.Sleep(1 * time.Second)
		log.Println("got:", msg)
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
	"time"

	"github.com/alextanhongpin/pkg/observable"
)

var (
	LoginEvent   = observable.NewEvent("login")
	LogoutEvent  = observable.NewEvent("logout")
	UnknownEvent = observable.NewEvent("unknown")
)

type UserService struct {
	observable.Observer
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
	err := u.Emit(LoginEvent, req)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	usersvc := &UserService{observable.NewAsync(10)}
	// Start n workers.
	for i := 0; i < 2; i++ {
		usersvc.Start()
	}

	// Logout event.
	usersvc.On(LogoutEvent, func(msg interface{}) error {
		fmt.Println("on logout:", msg)
		return nil
	})
	usersvc.Logout("john has logged out")

	// Register the login event.
	usersvc.On(LoginEvent, func(msg interface{}) error {
		// Fake work done here.
		time.Sleep(1 * time.Second)
		switch v := msg.(type) {
		case LoginRequest:
			fmt.Println("on login:", v.Msg)
		default:
			fmt.Println("not handled")
		}
		return nil
	})

	// Trigger the login event.
	go func() {
		for i := 0; i < 1000; i++ {
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			usersvc.Login(LoginRequest{fmt.Sprintf("user %d", i+1)})
		}
	}()

	// Debug goroutine count.
	// go func() {
	//         ticker := time.NewTicker(1 * time.Second)
	//         defer ticker.Stop()
	//         for {
	//                 select {
	//                 case <-ticker.C:
	//                         log.Println("goroutines:", runtime.NumGoroutine())
	//                 }
	//         }
	// }()

	time.Sleep(3 * time.Second)
	usersvc.Stop()
	log.Println("shutdown")
}
```
