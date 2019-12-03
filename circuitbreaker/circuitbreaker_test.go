package circuitbreaker_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/alextanhongpin/pkg/circuitbreaker"
)

func Example() {
	state := circuitbreaker.NewDefaultState()
	state.Timeout = 1 * time.Second
	cb := New(state)
	for i := 0; i < 10; i++ {
		res, err := cb.Handle(func() (interface{}, error) {
			return nil, errors.New("some error")
		})
		fmt.Println(res, err)
	}
	fmt.Println("sleep 1,1 seconds")
	time.Sleep(1100 * time.Millisecond)

	for i := 0; i < 3; i++ {
		res, err := cb.Handle(func() (interface{}, error) {
			return nil, errors.New("another error")
		})
		fmt.Println(res, err)
	}

	fmt.Println("sleep 1.1 seconds")
	time.Sleep(1100 * time.Millisecond)
	for i := 0; i < 15; i++ {
		res, err := cb.Handle(func() (interface{}, error) {
			return true, nil
		})
		fmt.Println(res, err)
	}
	for i := 0; i < 20; i++ {
		res, err := cb.Handle(func() (interface{}, error) {
			return nil, errors.New("some error")
		})
		fmt.Println(res, err)
	}
}
