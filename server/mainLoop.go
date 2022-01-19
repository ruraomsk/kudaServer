package server

import "github.com/ruraomsk/ag-server/logger"

func (d *deviceInfo) mainLoop() {
	logger.Info.Printf("Начинаем обмен с %d", d.uid)
	go d.ReadMessage()
	go d.WriteMessage()
	// defer func() {
	// 	// close(d.readChan)
	// 	// close(d.writeChan)
	// 	// close(d.command)
	// }()
	d.writeChan <- d.getMeStatus()

	for {
		select {
		case message, ok := <-d.readChan:
			if !ok {
				//Канал ввода закрылся
				return
			}
			logger.Debug.Printf("Пришло %v", message)
		}
	}
}
