package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/kudaServer/server"
	"github.com/ruraomsk/kudaServer/setup"
)

var (
	//go:embed config
	config embed.FS
)

func init() {
	setup.Set = new(setup.Setup)
	if _, err := toml.DecodeFS(config, "config/config.toml", &setup.Set); err != nil {
		fmt.Println("Отсутствует config.toml")
		os.Exit(-1)
		return
	}
	os.MkdirAll(setup.Set.LogPath, 0777)
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := logger.Init(setup.Set.LogPath); err != nil {
		log.Panic("Error logger system", err.Error())
		return
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	fmt.Println("kudaServer start")
	logger.Info.Println("kudaServer start")
	server.StartServer(2018)
	for {
		<-c
		fmt.Println("Wait make abort...")
		time.Sleep(3 * time.Second)
		fmt.Println("kudaServer stop")
		logger.Info.Println("kudaServer stop")
		os.Exit(0)
	}

}
