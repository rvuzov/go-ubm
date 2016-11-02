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

func (m *MetricPush) Receive() interface{} {
	ubm.Metrics.Push(m.UserID, m.Key, m.Value)
	return service.NewSuccess()
}
