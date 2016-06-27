package metric

import (
	"../../../bmodel"
	"../service"
)

type (
	MetricPush struct {
		UserID string `json:"userID`
		Key    string `json:"key"`
		Value  int    `json:"value"`
	}
)

func (msg *MetricPush) Receive() interface{} {

	err := bmodel.Metrics.Push(
		(*msg).UserID,
		(*msg).Key,
		(*msg).Value,
	)

	if err != nil {
		return service.NewError(err.Error())
	}

	return service.NewSuccess()
}
