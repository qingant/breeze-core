package main

import (
	"./breeze"
	"fmt"
	"runtime"
	"syscall"
)
var Build string
var Revision string
var Version string

func systemInit() {
	runtime.GOMAXPROCS(9);
	var rLimit syscall.Rlimit
	rLimit.Cur = 19999
	rLimit.Max = 19999
	syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}

func main() {
	// fmt.Println("hello world")
	systemInit()
	fmt.Printf("breeze-core build: %s\nrevision: %s\nversion: %s\n", Build, Revision, Version)
	breeze.InitApp()
	breeze.StartApp()
	// c := breeze.ReadConfig("./config.toml")
	// breeze.StartServer(c.Address)
	// breeze.TestStore()
	// breeze.TestStrategy()
}
