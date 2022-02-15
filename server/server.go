package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ruraomsk/ag-server/comm"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ag-server/pudge"
)

type Message struct {
	Messages map[string][]byte `json:"messages"`
}
type deviceInfo struct {
	socket    net.Conn
	uid       int
	key       []byte
	readChan  chan Message
	writeChan chan Message
	toutRead  time.Duration
	toutWrite time.Duration
	ctrl      *pudge.Controller
	cross     *pudge.Cross
	command   chan comm.CommandARM
	work      bool
}

var (
	messageNotFound = "notFound!"
	devs            struct {
		sync.Mutex
		devs map[int]deviceInfo
	}
)

func listenConnect(port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Error.Printf("Ошибка открытия порта %d", port)
		return
	}
	for {
		socket, err := ln.Accept()
		if err != nil {
			logger.Error.Printf("Ошибка accept %s", err.Error())

		}
		logger.Info.Printf("Подключаем %s", socket.RemoteAddr().String())
		go workerDevice(socket)
	}

}
func StartServer(port int) {
	devs.devs = make(map[int]deviceInfo)
	go listenConnect(port)
}
