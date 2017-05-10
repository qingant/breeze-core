package breeze

type IExecutor interface {
	GetType() string
	Call(uc *UserContext, action *Action, event *Event) (EventList, error)
}

type ExecutorManager struct {
	Executors map[string]IExecutor
}

var em = &ExecutorManager{Executors: make(map[string]IExecutor)}

func GetExecutorManager() *ExecutorManager {
	return em
}

// func NewExecutor(type_ string, host string, port int) *Executor {
// 	executor := &Executor{Type: type_, Host: host, Port: port}
// 	em.Executors[type_] = executor
// 	return executor
// }

func (em *ExecutorManager) GetExecutor(type_ string) IExecutor {
	return em.Executors[type_]
}

func (em *ExecutorManager) AddExecutor(type_ string, executor IExecutor) {
	em.Executors[executor.GetType()] = executor
}
