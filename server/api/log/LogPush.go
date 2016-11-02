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

func (m *LogPush) Receive() interface{} {
	err := ubm.Logs.Push(m.UserId, m.Key, m.Value)
	if err != nil {
		return service.NewError(err.Error())
	}
	return service.NewSuccess()
}
