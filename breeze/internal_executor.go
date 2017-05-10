package breeze

import (
	"github.com/kr/pretty"
)

type logExecutor struct {
}

func (e *logExecutor) GetType() string {
	return "log"
}

func (e *logExecutor) Call(uc *UserContext, action *Action, event *Event) (EventList, error) {
	pretty.Println("UserContext: ", *uc)
	pretty.Println("Action: ", *action)
	pretty.Println("Event: ", *event)
	return nil, nil
}
