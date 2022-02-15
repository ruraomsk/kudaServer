package server

import (
	"bufio"
	"encoding/base64"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/ruraomsk/ag-server/comm"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ag-server/pudge"
)

func isDeviceIntoCross(uid int) (*pudge.Cross, bool) {
	if uid == 0 {
		return nil, false
	}
	return new(pudge.Cross), true
}
func isDeviceIntoDevices(uid int) (*pudge.Controller, bool) {
	return new(pudge.Controller), true
}
func createDeviceToDevices(cross *pudge.Cross) *pudge.Controller {
	return new(pudge.Controller)
}
func (d *deviceInfo) setTouts() {
	d.ctrl.TMax = 400
	d.toutRead = time.Second * time.Duration(d.ctrl.TMax+60)
	d.toutWrite = time.Second * 10
	d.ctrl.TimeOut = int64(d.toutRead)
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
	number = strings.ReplaceAll(number, "\n", "")
	uid, err := strconv.Atoi(number)
	if err != nil {
		logger.Error.Printf("Неверный номер от %s %s", socket.RemoteAddr(), err.Error())
		return
	}
	var (
		cross *pudge.Cross      = &pudge.Cross{}
		ctrl  *pudge.Controller = &pudge.Controller{}
		is    bool
		dev   deviceInfo
	)
	if cross, is = isDeviceIntoCross(uid); !is {
		logger.Info.Printf("Устройство %d не подключено", uid)
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
		// close(dev.readChan)
		// close(dev.writeChan)
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
		work:      true,
	}
	dev.generateKey(16)
	dev.setTouts()
	dev.ctrl.ID = dev.uid
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
