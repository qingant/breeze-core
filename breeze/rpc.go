package breeze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kr/pretty"
	"log"
	"net/http"
	"time"
)

type Executor struct {
	Type string
	Host string
	Port int
}

func (e *Executor) GetType() string {
	return e.Type
}

func (e *Executor) Call(uc *UserContext, action *Action, event *Event) (EventList, error) {
	defer timeTrack(time.Now(), fmt.Sprintf("%s:%s.%s.%s", uc.Strategy, action.Module, action.Class, action.Func))
	url := fmt.Sprintf("http://%s:%d/v1/call", e.Host, e.Port)
	log.Printf("URL: %s", url)
	stra := GetStrategyManager().GetStrategy(uc.Strategy)
	params := map[string]interface{}{
		"user_context": uc,
		"action":       action,
		"event":        event,
		"strategy":     stra,
		"addr":         GetConfig().Address,
		"log":          GetConfig().Redis,
		"db":           GetConfig().Db}
	buf, _ := json.Marshal(params)
	ret, err := http.Post(url, "application/json", bytes.NewReader(buf))
	if err != nil {
		GetStatus().logError(err.Error())
		return nil, err
	}
	decoder := json.NewDecoder(ret.Body)
	el := &struct {
		Status string
		Events EventList
	}{}
	decoder.Decode(el)
	// json, err := sjson.NewJson(body)
	pretty.Println("EVents: ", el)
	return el.Events, nil
}
