package grace

import (
	"context"
	"sync"
)

type Shutdown func(context.Context)

type Shutdowns []Shutdown

func NewShutdowns() Shutdowns {
	return make([]Shutdown, 0)
}

func (shutdowns *Shutdowns) Append(shutdown Shutdown) {
	*shutdowns = append(*shutdowns, shutdown)
}

func (shutdowns Shutdowns) Close(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(shutdowns))
	for _, shutdown := range shutdowns {
		go func(shutdown Shutdown) {
			defer wg.Done()
			shutdown(ctx)
		}(shutdown)
	}
	wg.Wait()
}
