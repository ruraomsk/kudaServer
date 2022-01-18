package server

import (
	"bufio"
	"encoding/base64"
	"net"
	"strconv"
	"time"

	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/ag-server/comm"
	"github.com/ruraomsk/ag-server/pudge"
)

func isDeviceIntoCross(uid int) (*pudge.Cross, bool) {
	return new(pudge.Cross), true
}
func isDeviceIntoDevices(uid int) (*pudge.Controller, bool) {
	return new(pudge.Controller), true
}
func createDeviceToDevices(cross *pudge.Cross) *pudge.Controller {
	return new(pudge.Controller)
}
func (d *deviceInfo) setTouts() {
	d.toutRead = time.Second * 10
	d.toutWrite = time.Second * 10
}
func workerDevice(socket net.Conn) {
	defer socket.Close()
	socket.SetDeadline(time.Now().Add(time.Second * 10))
	reader := bufio.NewReader(socket)
	writer := bufio.NewWriter(socket)
	// Вначале считываем номер устройства
	number, err := reader.ReadString('\n')
	if err != nil {
		logger.Error.Printf("Чтение номера от %s %s", socket.RemoteAddr(), err.Error())
		return
	}
	uid, err := strconv.Atoi(number)
	if err != nil {
		logger.Error.Printf("Неверный номер от %s %s", socket.RemoteAddr(), err.Error())
		return
	}
	var (
		cross *pudge.Cross
		ctrl  *pudge.Controller
		is    bool
		dev   deviceInfo
	)
	if cross, is = isDeviceIntoCross(uid); !is {
		writer.WriteString(messageNotFound)
		writer.WriteString("\n")
		_ = writer.Flush()
		return
	}
	if ctrl, is = isDeviceIntoDevices(uid); !is {
		ctrl = createDeviceToDevices(cross)
	}

	devs.Lock()

	dev, is = devs.devs[uid]
	if is {
		dev.socket.Close()
		close(dev.readChan)
		close(dev.writeChan)
	}
	dev = deviceInfo{
		socket:    socket,
		uid:       uid,
		readChan:  make(chan Message),
		writeChan: make(chan Message),
		toutRead:  0,
		toutWrite: 0,
		ctrl:      ctrl,
		cross:     cross,
		command:   make(chan comm.CommandARM),
	}
	dev.generateKey(16)
	dev.setTouts()
	devs.devs[uid] = dev
	writer.WriteString(base64.StdEncoding.EncodeToString(dev.key))
	writer.WriteString("\n")

	devs.Unlock()

	err = writer.Flush()

	if err != nil {
		logger.Error.Printf("Устройство %d передача сообщения для %s %s", uid, socket.RemoteAddr(), err.Error())
		devs.Lock()
		delete(devs.devs, uid)
		devs.Unlock()
		return
	}
	dev.mainLoop()
}
