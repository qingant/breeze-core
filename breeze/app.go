package breeze

import (
	"github.com/kr/pretty"
	"os"
	//"../redis"
)

//var r *redis.Server
//
//func GetRedis() *redis.Server {
//	return r
//}
//func InitRedis() {
//	address := fmt.Sprintf("%s:%d", GetConfig().Redis.Host, GetConfig().Redis.Port)
//	r = redis.NewServer(address)
//	go r.Serve()
//}


func InitApp() {
	var path string
	if os.Getenv("ENV") == "prod" {
		path = "./config.prod.toml"
	} else if os.Getenv("ENV") == "dev" {
		path = "./config.dev.toml"
	} else if os.Getenv("ENV") == "staging" {
		path = "./config.staging.toml"
	} else {
		path = "./config.toml"
	}
	c := ReadConfig(path)
	// InitRedis()
	InitStore(c.Db.Host, c.Db.Port)
	InitStatus()
	InitStrategies()
	InitContexts()
	GetScheduler().Init(c.Tasks)
	if _, err := os.Stat(c.StraPath); os.IsNotExist(err) {
		os.Mkdir(c.StraPath, os.FileMode(0666))
	}
	pretty.Println("Config: ", *c)
	l := &logExecutor{}

	GetExecutorManager().AddExecutor(l.GetType(), l)
}

func StartApp() {

	StartServer(GetConfig().Address)
}
