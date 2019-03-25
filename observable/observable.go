package observable

// Observer represents the observer interface.
type Observer interface {
	On(event Event, fn Action)
	Emit(event Event, params interface{}) error
	Start()
	Stop()
}

// Event represents the type of Event that is registered.
type Event string

func (e Event) String() string {
	return string(e)
}

// NewEvent returns a new Event.
func NewEvent(s string) Event {
	return Event(s)
}

// Action represents the function that is executed for the given Event.
type Action func(interface{}) error

// Message represents the payload that is sent through the channel.
type Message struct {
	event  Event
	params interface{}
}
