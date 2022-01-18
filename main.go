package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/ruraomsk/TLServer/logger"
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

}
