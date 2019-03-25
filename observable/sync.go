package observable

import (
	"fmt"
	"sync"
)

type SyncObservable struct {
	sync.RWMutex
	events map[Event][]Action
}

func NewSync() *SyncObservable {
	events := make(map[Event][]Action)
	return &SyncObservable{events: events}
}

func (o *SyncObservable) Start() {}
func (o *SyncObservable) Stop()  {}

func (o *SyncObservable) On(event Event, fn Action) {
	o.Lock()
	_, exist := o.events[event]
	if !exist {
		o.events[event] = make([]Action, 0)
	}
	o.events[event] = append(o.events[event], fn)
	o.Unlock()
}

func (o *SyncObservable) Emit(event Event, params interface{}) error {
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
