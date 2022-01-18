package server

func (d *deviceInfo) mainLoop() {
	go d.ReadMessage()
	go d.WriteMessage()
	defer func() {
		close(d.readChan)
		close(d.writeChan)
		close(d.command)
	}()
	d.writeChan <- d.getMeStatus()
	for {
		select {
		case message, ok := <-d.readChan:
			if !ok {
				//Канал ввода закрылся
				return
			}

		}
	}
}
