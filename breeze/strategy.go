package breeze

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/kr/pretty"
	"io/ioutil"
	// "os"
	"sync"
	"time"
	"gopkg.in/mgo.v2/bson"
)

type EventMap map[string]ActionList
type ParamsType map[string]interface{}
type DependencyList []string

// Strategy is a execution map
type Strategy struct {
	ID           string
	Route        EventMap
	Params       interface{}
	Tag          string
	Labels       interface{}
	Title        string
	Desc         interface{}
	Slogan       string
	Visiblity    bool
	Packageless  bool
	Dependencies DependencyList
}

type StrategyPkg struct {
	ID      string
	Update  time.Time
	Content []byte
}
type StrategyManager struct {
	strategies map[string]*Strategy
	lock       sync.RWMutex
}

var sm = &StrategyManager{strategies: make(map[string]*Strategy)}

func GetStrategyManager() *StrategyManager {
	return sm
}
func (sm *StrategyManager) AddStrategy(s *Strategy) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	pretty.Println("Add Strategy: ", s)
	if s == nil {
		return
	}
	sm.strategies[s.ID] = s
	GetStore().InsertDoc(STRATEGY_TBL_NAME, *s, s.ID)
}

func (sm *StrategyManager) GetStrategy(id string) *Strategy {
	return sm.strategies[id]
}

func (sm *StrategyManager) Remove(id string) int {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	_, exists := sm.strategies[id]
	if exists {
		delete(sm.strategies, id)
		return GetStore().Remove(STRATEGY_TBL_NAME, id)
	}
	return 0
}

func (s *Store) ReadStrategy(id string) *Strategy {
	strategy := &Strategy{}
	s.ReadObject(STRATEGY_TBL_NAME, id, strategy)
	GetStrategyManager().AddStrategy(strategy)
	return strategy
}

func NewStrategyFromToml(data string) (*Strategy, error) {
	s := &Strategy{}
	_, err := toml.Decode(data, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Strategy) path() string {
	path := fmt.Sprintf("%s/%s_%s.zip", GetConfig().StraPath, s.ID, s.Tag)
	fmt.Println("Write: ", path)
	return path
}
func (s *Strategy) name() string {
	name := fmt.Sprintf("%s_%s", s.ID, s.Tag)
	return name
}
func DeployStrategy(content []byte) (*Strategy, error) {
	r := bytes.NewReader(content)
	stra, err := zip.NewReader(r, int64(len(content)))
	if err != nil {
		return nil, err
	}
	for _, file := range stra.File {
		fmt.Println("File: ", file.Name)
		if file.Name == "stra.toml" {
			r, _ := file.Open()
			buf, _ := ioutil.ReadAll(r)
			defer r.Close()
			strategy, err := NewStrategyFromToml(string(buf))
			if err != nil {
				return nil, err
			}

			// ioutil.WriteFile(strategy.path(), content, os.FileMode(0666))
			GetStore().InsertDoc(STRATEGY_PKG_TBL_NAME, StrategyPkg{
				ID:      strategy.name(),
				Update:  time.Now(),
				Content: content}, strategy.name())
			GetStrategyManager().AddStrategy(strategy)
			return strategy, nil
		}
	}
	return nil, errors.New("Not a strategy")

}

func InitStrategies() {
	var ss []*Strategy
	GetStore().ReadAll(STRATEGY_TBL_NAME, &ss)
	pretty.Println(ss)
	for _, s := range ss {
		GetStrategyManager().strategies[s.ID] = s
		s.Params = bsonToObject(s.Params.(bson.M))
	}
}
