package breeze

// Action is a slots set
type Action struct {
	Type   string
	Module string
	Class  string
	Func   string
	Params map[string]interface{}
}

// ActionList : list of Component
type ActionList []Action
