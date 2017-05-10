package breeze

import "github.com/BurntSushi/toml"
import "fmt"

type dbConfig struct {
	Host string
	Port int
}

type redisConfig struct {
	Host string
	Port int
}

// Config for breeze-core
type Config struct {
	Address   string
	Db        dbConfig
	Redis     redisConfig
	StraPath  string
	Executors map[string]Executor
	Tasks     map[string]ConfigTask
}

var config *Config = &Config{}

func ReadConfig(path string) *Config {
	if _, err := toml.DecodeFile(path, config); err != nil {
		panic(err)
	}
	for k, v := range config.Executors {
		fmt.Printf("Executor <%s>: %v", k, v)
		GetExecutorManager().AddExecutor(v.GetType(), &v)
	}
	return config
}

func GetConfig() *Config {
	return config
}
