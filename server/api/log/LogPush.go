package log

import (
	"github.com/Lamzin/go-ubm/ubm"
	"github.com/Lamzin/go-ubm/server/api/service"
)

type (
	LogPush struct {
		UserId string      `json:"userID"`
		Key    string      `json:"key"`
		Value  interface{} `json:"value"`
	}
)

func (msg *LogPush) Receive() interface{} {

	err := ubm.Logs.Push(
		(*msg).UserId,
		(*msg).Key,
		(*msg).Value,
	)

	if err != nil {
		return service.NewError(err.Error())
	}

	return service.NewSuccess()
}
