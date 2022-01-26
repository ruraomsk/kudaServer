package server

import (
	"encoding/json"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ag-server/pudge"
)

type BaseCtrl struct {
	ID         int       `json:"id"`
	TimeDevice time.Time `json:"dtime"`    // Время устройства
	TechMode   int       `json:"techmode"` //Технологический режим
	/*
		Технологический режим
		1 - выбор ПК по времени по суточной карте ВР-СК;
		2 - выбор ПК по недельной карте ВР-НК;
		3 - выбор ПК по времени по суточной карте, назначенной
		оператором ДУ-СК;
		4 - выбор ПК по недельной карте, назначенной оператором
		ДУ-НК;
		5 - план по запросу оператора ДУ-ПК;
		6 - резервный план (отсутствие точного времени) РП;
		7 – коррекция привязки с ИП;
		8 – коррекция привязки с сервера;
		9 – выбор ПК по годовой карте;
		10 – выбор ПК по ХТ;
		11 – выбор ПК по картограмме;
		12 – противозаторовое управление.
	*/
	Base    bool  `json:"base"` //Если истина то работает по базовой привязке
	PK      int   `json:"pk"`   //Номер плана координации
	CK      int   `json:"ck"`   //Номер суточной карты
	NK      int   `json:"nk"`   //Номер недельной карты
	TMax    int64 `json:"tmax"` //Максимальное время ожидания ответа от сервера в секундах
	TimeOut int64 `json:"tout"` //TimeOut на чтение от контроллера в секундах
}

var execute = map[string]interface{}{
	"base":            BaseCtrl{},
	"traffic":         pudge.Traffic{},
	"Status":          pudge.Status{},
	"StatusCommandDU": pudge.StatusCommandDU{},
	"DK":              pudge.DK{},
	"Model":           pudge.Model{},
	"ErrorDevice":     pudge.ErrorDevice{},
	"GPS":             pudge.GPS{},
	"Input":           pudge.Input{},
}

func (d *deviceInfo) mainLoop() {
	logger.Info.Printf("Начинаем обмен с %d", d.uid)
	go d.ReadMessage()
	go d.WriteMessage()
	go d.sendCommandTest()
	defer func() {
		close(d.command)
	}()
	d.writeChan <- d.getMeStatus()
	sendGetStatus := time.NewTicker(time.Duration(d.ctrl.TMax-20) * time.Second)
	deviceNotWork := time.NewTicker(time.Duration(d.ctrl.TMax+60) * time.Second)
	for {
		select {
		case message, ok := <-d.readChan:
			if !ok {
				//Канал ввода закрылся
				return
			}
			replay, need := d.updateDevice(message)
			if need {
				d.writeChan <- replay
			}
			deviceNotWork.Reset(time.Duration(d.ctrl.TMax+60) * time.Second)
		case <-sendGetStatus.C:
			d.writeChan <- d.getMeStatus()
		case <-deviceNotWork.C:
			logger.Error.Printf("Устройство %d не отвечает", d.uid)
			return
		case cmd := <-d.command:
			m := newMessage()
			body, _ := json.Marshal(cmd)
			m.Messages["command"] = body
			d.writeChan <- m
		}
	}
}
func (d *deviceInfo) updateDevice(m Message) (Message, bool) {
	need := false
	for name, buffer := range m.Messages {
		_, is := execute[name]
		if !is {
			logger.Error.Printf("Неизвестное %s", name)
			continue
		}
		// logger.Info.Printf("Обработка %s", name)
		// json.Unmarshal(buffer, &value)
		switch name {
		case "base":
			var v pudge.Controller
			json.Unmarshal(buffer, &v)
			d.ctrl.TechMode = v.TechMode
			d.ctrl.TimeDevice = v.TimeDevice
			d.ctrl.Base = v.Base
			d.ctrl.PK = v.PK
			d.ctrl.CK = v.CK
			d.ctrl.NK = v.NK
			d.ctrl.TMax = v.TMax
			d.ctrl.TimeOut = v.TimeOut
		case "traffic":
			var v pudge.Traffic
			json.Unmarshal(buffer, &v)
			d.ctrl.Traffic = v
		case "Status":
			var v pudge.Status
			json.Unmarshal(buffer, &v)
			d.ctrl.Status = v
		case "StatusCommandDU":
			var v pudge.StatusCommandDU
			json.Unmarshal(buffer, &v)
			d.ctrl.StatusCommandDU = v
		case "DK":
			var v pudge.DK
			json.Unmarshal(buffer, &v)
			d.ctrl.DK = v
		case "Model":
			var v pudge.Model
			json.Unmarshal(buffer, &v)
			d.ctrl.Model = v
		case "ErrorDevice":
			var v pudge.ErrorDevice
			json.Unmarshal(buffer, &v)
			d.ctrl.Error = v
		case "GPS":
			var v pudge.GPS
			json.Unmarshal(buffer, &v)
			d.ctrl.GPS = v
		case "Input":
			var v pudge.Input
			json.Unmarshal(buffer, &v)
			d.ctrl.Input = v
		}
	}
	logger.Info.Printf("ctrl %v", d.ctrl)
	return Message{}, need
}
