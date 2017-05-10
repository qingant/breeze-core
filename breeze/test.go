package breeze

import (
	"fmt"
	"github.com/kr/pretty"
	"io/ioutil"
	"log"
)

func TestUc() {
	c := ReadConfig("./config.toml")
	log.Println(c)
	uc, _ := NewUserContext("abcd", "qingant", "demo")
	executor := &Executor{"python", "localhost", 5000}
	action := &Action{Type: "python", Module: "breeze", Class: "Test", Func: "test"}
	event := &Event{Type: "dayline", Params: []interface{}{"a", "b"}}
	executor.Call(uc, action, event)
}

func TestStore() {
	c := ReadConfig("./config.toml")
	InitStore(c.Db.Host, c.Db.Port)

	s := &Strategy{ID: "test", Route: EventMap{"dayline": ActionList{Action{Type: "python",
		Module: "test", Class: "Test", Func: "on_dayline"}}}, Params: ParamsType{}}
	GetStore().InsertDoc("strategies", s, s.ID)
	ss := &Strategy{}
	GetStore().ReadObject("strategies", "test", ss)
	fmt.Println("S: ", s)
	fmt.Println("SS: ", ss)

}

func TestStrategy() {
	buf, _ := ioutil.ReadFile("./examples/simple_strategy.toml")
	s, _ := NewStrategyFromToml(string(buf))
	// fmt.Println("S: ", s)
	pretty.Println("S: ", s)
}
