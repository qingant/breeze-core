package breeze

import (
	"github.com/go-redis/redis"
	"fmt"
	"encoding/json"
	"time"
)

type Status struct {
	redis *redis.Client
}

var s *Status = nil

func GetStatus() *Status {
	return s
}

func InitStatus() {
	addr := fmt.Sprintf("%s:%d", GetConfig().Redis.Host, GetConfig().Redis.Port)
	s = &Status{
		redis.NewClient(&redis.Options{
			Addr: addr,
			DB: 1,
		}),
	}
}

func (s *Status) logEvent(ucid string, e *Event) {
	channel := fmt.Sprintf("event.%s", ucid)
	_e := &struct {
		Time time.Time
		Ucid   string
		Event *Event
	}{time.Now(), ucid, e}
	GetStore().InsertDoc(EVENT_TBL_NAME, _e, "")
	buf, _ := json.Marshal(_e)
	repr := string(buf)
	s.redis.Publish(channel, repr)
	s.redis.Publish("events", repr)
}

func (s *Status) logPerf(name string, elapsed time.Duration) {
	desc := fmt.Sprintf("%s: %s", name, elapsed)
	s.redis.Publish("perfs", desc)
	// s.redis.Publish(fmt.Sprintf("perf.%s", name), string(elapsed))
}

func (s *Status) logError(desc string) {
	s.redis.Publish("errors", desc)
}