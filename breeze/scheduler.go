package breeze

import (
	"github.com/robfig/cron"
	"log"
)

type Scheduler struct {
	Cron *cron.Cron
}

type ConfigTask struct {
	Spec  string
	Event string
}

var scheduler = &Scheduler{cron.New()}

func GetScheduler() *Scheduler {
	return scheduler
}

func (s *Scheduler) Add(spec string, ucid string, event *Event) {
	s.Cron.AddFunc(spec, func() {
		if ucid == "" {
			GetContextManager().Broadcast(event)
		} else {
			GetContextManager().SendEvent(ucid, event)
		}

	})
}

func (s *Scheduler) Init(tasks map[string]ConfigTask) {
	for _, v := range tasks {
		event := &Event{Type: v.Event}
		log.Println("Add Task:", tasks)
		s.Add(v.Spec, "", event)
	}
}
