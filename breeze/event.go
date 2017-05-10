package breeze

import "time"

// Event
type Event struct {
	Type   string
	Params interface{}
	From   string
	Delay  time.Duration
	callback func (interface{})
}

var StopEvent = &Event{}

type EventList []*Event
