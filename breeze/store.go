package breeze

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"log"
	"gopkg.in/mgo.v2/bson"
	"github.com/kr/pretty"
)

const BREEZE_DB = "breeze"
const STRATEGY_TBL_NAME = "strategies"
const STRATEGY_PKG_TBL_NAME = "strategy_pkgs"
const USER_CONTEXT_TBL_NAME = "user_context"
const EVENT_TBL_NAME = "event"
const TRADE_TBL_NAME = "trade"

type Store struct {
	Host    string
	Port    int
	Db      string
	session *mgo.Session
}

var store *Store = nil

func InitStore(host string, port int) {
	store = &Store{Host: host, Port: port, Db: BREEZE_DB}
	store.connect()
}

func GetStore() *Store {
	return store
}

func (s *Store) connect() error {
	url := fmt.Sprintf("%s:%d", s.Host, s.Port)
	log.Println("Url: ", url)
	var err error
	s.session, err = mgo.Dial(url)
	return err
}

func (s *Store) InsertDoc(collection string, doc interface{}, id string) error {
	log.Printf("Insert into `%s`: %v", collection, doc)
	db := s.session.DB(s.Db)
	fmt.Println("Db: ", db)
	if id != "" {
		n, err := db.C(collection).RemoveAll(map[string]string{"id": id})
		fmt.Println("Error:", n, err)
	}
	return db.C(collection).Insert(doc)
}

func (s *Store) ReadObject(collection string, id string, obj interface{}) {
	c := s.session.DB(s.Db).C(collection)
	c.Find(bson.M{"id": id}).One(obj)
	log.Println(obj)
}

func (s *Store) ReadAll(collection string, obj interface{}) {
	c := s.session.DB(s.Db).C(collection)
	c.Find(nil).All(obj)
}

func (s *Store) ReadBy(collection string, obj interface{}, condition interface{}) {
	c := s.session.DB(s.Db).C(collection)
	c.Find(condition).All(obj)
}

func (s *Store) Remove(collection string, id string) int {
	c := s.session.DB(s.Db).C(collection)
	n, _ := c.RemoveAll(map[string]string{"id": id})
	if n != nil {
		return n.Removed
	} else {
		return 0
	}
}

func bsonToObject(m bson.M) interface{} {
	o := make(map[string]interface{})
	pretty.Println("M:", m)
	// buf, _ := bson.MarshalJSON(m)
	for k, e := range m {
		fmt.Println("B:", k, e)
		switch e.(type) {
		case int64:
			o[k] = e
		case float64:
			o[k] = e
		case string:
			o[k] = e
		case bool:
			o[k] = e
		case bson.M:
			o[k] = bsonToObject(e.(bson.M))
		default:
			o[k] = e

		}
	}
	pretty.Println("O:", o)
	return o
	// pretty.Println("J:", string(buf))
	// var o interface{}
	// json.Unmarshal(buf, &o)
	pretty.Println("O:", o)
	return o
}
