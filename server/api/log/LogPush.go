package log

import (
	"../../../bmodel"
	"../service"
)

type (
	LogPush struct {
		UserId string      `json:"userID"`
		Key    string      `json:"key"`
		Value  interface{} `json:"value"`
	}
)

func (msg *LogPush) Receive() interface{} {

	err := bmodel.Logs.Push(
		(*msg).UserId,
		(*msg).Key,
		(*msg).Value,
	)

	if err != nil {
		return service.NewError(err.Error())
	}

	return service.NewSuccess()
}
