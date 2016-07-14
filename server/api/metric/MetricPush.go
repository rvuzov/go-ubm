package metric

import (
	"github.com/Lamzin/go-ubm/ubm"
	"github.com/Lamzin/go-ubm/server/api/service"
)

type (
	MetricPush struct {
		UserID string `json:"userID`
		Key    string `json:"key"`
		Value  int    `json:"value"`
	}
)

func (msg *MetricPush) Receive() interface{} {

	err := ubm.Metrics.Push(
		(*msg).UserID,
		(*msg).Key,
		(*msg).Value,
	)

	if err != nil {
		return service.NewError(err.Error())
	}

	return service.NewSuccess()
}
