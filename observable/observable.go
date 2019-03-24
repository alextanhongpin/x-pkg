package observable

import (
	"fmt"
	"sync"
)

type ObserverFunc func(interface{}) error
type Observer interface {
	On(event string, fn ObserverFunc)
	Emit(event string, params interface{}) error
}

type SyncObservable struct {
	sync.RWMutex
	events map[string][]ObserverFunc
}

func (o *SyncObservable) On(event string, fn ObserverFunc) {
	o.Lock()
	_, exist := o.events[event]
	if !exist {
		o.events[event] = make([]ObserverFunc, 0)
	}
	o.events[event] = append(o.events[event], fn)
	o.Unlock()
}

func (o *SyncObservable) Emit(event string, params interface{}) error {
	o.RLock()
	fns, exist := o.events[event]
	o.RUnlock()
	if !exist {
		return fmt.Errorf(`event "%s" does not exist`, event)
	}
	for _, fn := range fns {
		err := fn(params)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewSyncObservable() *SyncObservable {
	events := make(map[string][]ObserverFunc)
	return &SyncObservable{events: events}
}
