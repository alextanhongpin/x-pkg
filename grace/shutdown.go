package grace

import (
	"context"
	"sync"
)

type Shutdown func(context.Context)

type ShutdownGroup []Shutdown

func NewShutdownGroup() ShutdownGroup {
	return make(ShutdownGroup, 0)
}

func (sg *ShutdownGroup) Add(shutdown Shutdown) {
	*sg = append(*sg, shutdown)
}

func (sg ShutdownGroup) Close(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(sg))
	for _, s := range sg {
		go func(shutdown Shutdown) {
			defer wg.Done()
			shutdown(ctx)
		}(s)
	}
	wg.Wait()
}
