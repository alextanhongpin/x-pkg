package grace

import (
	"context"
	"sync"
)

type Shutdown func(context.Context)

type ShutdownGroup struct {
	s    []Shutdown
	once sync.Once
}

// NewShutdownGroup returns a new pointer to ShutdownGroup.
func NewShutdownGroup() *ShutdownGroup {
	return &ShutdownGroup{}
}

// Add a new Shutdown function to the list - each shutdown accepts a context to allow cancellation/timeout.
func (sg *ShutdownGroup) Add(shutdown Shutdown) {
	sg.s = append(sg.s, shutdown)
}

// Close will synchronize the termination of all resources, and will only execute once.
func (sg *ShutdownGroup) Close(ctx context.Context) {
	sg.once.Do(func() {
		var wg sync.WaitGroup
		wg.Add(len(sg.s))
		for _, s := range sg.s {
			go func(shutdown Shutdown) {
				defer wg.Done()
				shutdown(ctx)
			}(s)
		}
		wg.Wait()
	})
}
