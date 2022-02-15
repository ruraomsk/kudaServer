package server

import (
	"time"

	"github.com/ruraomsk/ag-server/comm"
)

var execTests = []comm.CommandARM{
	{Command: 2, Params: 1},
	{Command: 2, Params: 0},
	{Command: 4, Params: 1},
	{Command: 5, Params: 3},
	{Command: 6, Params: 2},
	{Command: 7, Params: 1},
	{Command: 9, Params: 9},
	{Command: 12, Params: 0},
	{Command: 4, Params: 0},
	{Command: 5, Params: 0},
	{Command: 6, Params: 0},
	{Command: 7, Params: 0},
	{Command: 9, Params: 0},
}

func (d *deviceInfo) sendCommandTest() {
	for {
		for _, c := range execTests {
			if d.work {
				d.command <- c
				time.Sleep(5 * time.Second)
			} else {
				return
			}
		}
	}
}
