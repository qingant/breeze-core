package breeze

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"sync"
	"github.com/kr/pretty"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"log"
)

const PipeDepth = 1024

type ContextStatus int

const (
	CREATED ContextStatus = iota
	READY
	STOPPED
)

type UserContext struct {
	ID       string
	UID      string
	Params   interface{}
	Strategy string
	pipe     chan *Event `bson:"-"`
	status   ContextStatus
}

type ContextManager struct {
	contexts map[string]*UserContext
	mutex    sync.RWMutex
}

var cm = &ContextManager{contexts: make(map[string]*UserContext)}

func GetContextManager() *ContextManager {
	return cm
}
func (cm *ContextManager) AddUC(uc *UserContext) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if context, _ := cm.contexts[uc.ID]; context != nil {
		context.Send(StopEvent)
	}
	cm.contexts[uc.ID] = uc
	GetStore().InsertDoc(USER_CONTEXT_TBL_NAME, *uc, uc.ID)
	uc.Start()
}

func (cm *ContextManager) GetUC(id string) *UserContext {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.contexts[id]
}

func (cm *ContextManager) SendEvent(id string, event *Event) error {
	fmt.Println("SendEvent: ", id, event)
	uc := cm.GetUC(id)
	if uc == nil {
		return errors.New("Context not exists")
	}
	uc.Send(event)
	return nil
}

func (cm *ContextManager) Broadcast(event *Event) (int, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	count := 0
	for k, v := range cm.contexts {
		fmt.Printf("Send %s %v\n", k, event)
		v.Send(event)
		count++
	}
	return count, nil
}

func NewUserContext(id string, uid string, strategy string) (*UserContext, error) {
	uc := &UserContext{ID: id,
		UID:           uid,
		Strategy:      strategy,
		pipe:          make(chan *Event, PipeDepth)}
	cm.AddUC(uc)
	return uc, nil
}

func NewUserContextFromJson(decoder *json.Decoder) (*UserContext, error) {
	uc := &UserContext{}
	err := decoder.Decode(uc)
	uc.pipe = make(chan *Event, PipeDepth)
	if uc.ID == "" {
		uc.ID = uuid.NewV4().String()
	}
	if err != nil {
		return nil, err
	}
	return uc, nil

}

func (s *Store) ReadUserContext(id string) *UserContext {
	uc := &UserContext{}
	s.ReadObject(USER_CONTEXT_TBL_NAME, id, uc)
	uc.pipe = make(chan *Event, PipeDepth)
	cm.AddUC(uc)
	return uc
}

func (uc *UserContext) getActionList(stra *Strategy, event *Event) ActionList {
	al, exists := stra.Route[event.Type]
	if exists {
		return al
	} else {
		for _, s := range stra.Dependencies {
			sid := strings.Split(s, ":")[0]
			stra := GetStrategyManager().GetStrategy(sid)
			al := uc.getActionList(stra, event)
			if al != nil {
				return al
			}
		}
		return nil
	}
}

func (uc *UserContext) onEvent(event *Event) {
	//GetStore().InsertDoc(EVENT_TBL_NAME, struct {
	//	Ucid string
	//	Time time.Time
	//	Event interface{}
	//}{uc.ID, time.Now(), event}, "")
	GetStatus().logEvent(uc.ID, event)
	strategy := GetStrategyManager().GetStrategy(uc.Strategy)
	if uc.status == STOPPED {
		return
	}
	// fmt.Println("Strategy: ", uc.Strategy, strategy)
	if strategy == nil {
		return
	}
	actionList := uc.getActionList(strategy, event)
	pretty.Println("ActionList:", actionList)
	if actionList == nil {
		return
	}
	for i := range actionList {
		action := actionList[i]
		fmt.Println(i, action)
		eventList, err := GetExecutorManager().GetExecutor(action.Type).Call(uc, &action, event)
		if event.callback != nil {
			event.callback(eventList)
			return
		}
		if err != nil {
			log.Println("RPC Error: ", err.Error())
			return
		}
		for j := range eventList {
			uc.onEvent(eventList[j])
			// uc.pipe <- eventList[j]
		}
	}
}

func (uc *UserContext) Send(event *Event) {
	fmt.Println("Send event")
	pretty.Println(event)
	uc.pipe <- event
}
func (uc *UserContext) stop() {
	uc.status = STOPPED
	uc.pipe <- StopEvent
}

func (uc *UserContext) loop() {
	// c := fmt.Sprintf("ready.%s", uc.ID)
	// m := fmt.Sprintf("%s ready", uc.ID)
	for {
		fmt.Println("loop", )
		fmt.Println("LLLen:", uc.ID, len(uc.pipe))
		event := <-uc.pipe
		if event == StopEvent {
			return
		}
		uc.onEvent(event)
	}
}

func (uc *UserContext) Start() {
	uc.Send(&Event{Type: "start", From: "builtin", Params: uc.Params})
	go uc.loop()
}

func InitContexts() {
	var ucs []*UserContext
	GetStore().ReadAll(USER_CONTEXT_TBL_NAME, &ucs)
	for _, uc := range ucs {
		fmt.Println("Init:", uc.ID)
		uc.pipe = make(chan *Event, PipeDepth)
		GetContextManager().contexts[uc.ID] = uc
		uc.Params = bsonToObject(uc.Params.(bson.M))
		uc.Start()
	}
}
