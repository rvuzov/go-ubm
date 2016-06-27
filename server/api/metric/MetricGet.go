package metric

import (
	"../../../bmodel"
	"../service"
)

type (
	MetricGet struct {
		UserID  string   `json:"userID"`
		Metrics []string `json:"metrics"`
	}

	MetricGetResponse struct {
		UserID  string         `json:"userID"`
		Metrics map[string]int `json:"metrics"`
	}
)

func (msg *MetricGet) Receive() interface{} {

	metrics, err := bmodel.Metrics.Get(
		(*msg).UserID,
		(*msg).Metrics,
	)

	if err != nil {
		return service.NewError(err.Error())
	}

	return MetricGetResponse{
		UserID:  (*msg).UserID,
		Metrics: metrics,
	}

}
