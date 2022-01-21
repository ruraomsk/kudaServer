package server

import (
	"bufio"
	"encoding/json"
	"strings"
	"time"

	"github.com/ruraomsk/ag-server/logger"
)

func (d *deviceInfo) ReadMessage() {
	defer d.socket.Close()
	defer close(d.readChan)
	reader := bufio.NewReader(d.socket)
	for {
		d.socket.SetReadDeadline(time.Now().Add(d.toutRead))
		message, err := reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("Устройство %d чтение сообщения от %s %s", d.uid, d.socket.RemoteAddr().String(), err.Error())
			return
		}
		message = strings.ReplaceAll(message, "\n", "")
		mess, err := d.decode(message)
		if err != nil {
			logger.Error.Printf("Устройство %d декодирование сообщения от %s %s", d.uid, d.socket.RemoteAddr().String(), err.Error())
			return
		}
		var inm Message
		err = json.Unmarshal(mess, &inm)
		if err != nil {
			logger.Error.Printf("Устройство %d unmarshal  сообщения %v %s", d.uid, message, err.Error())
			return
		}
		d.readChan <- inm
	}
}
func (d *deviceInfo) WriteMessage() {
	defer d.socket.Close()
	defer close(d.writeChan)
	writer := bufio.NewWriter(d.socket)
	for {
		message, ok := <-d.writeChan
		if !ok {
			logger.Error.Printf("канал вывода для %d закрыт", d.uid)
			return
		}
		buffer, err := json.Marshal(message)
		if err != nil {
			logger.Error.Printf("Устройство %d marshal  сообщения %v %s", d.uid, message, err.Error())
			return
		}
		d.socket.SetWriteDeadline(time.Now().Add(d.toutWrite))
		str, err := d.code(buffer)
		if err != nil {
			logger.Error.Printf("Устройство %d кодирование сообщения %v %s", d.uid, buffer, err.Error())
			return
		}
		_, _ = writer.WriteString(str)
		_, _ = writer.WriteString("\n")
		err = writer.Flush()
		if err != nil {
			logger.Error.Printf("Устройство %d передача сообщения для %s %s", d.uid, d.socket.RemoteAddr().String(), err.Error())
			return
		}
	}
}
