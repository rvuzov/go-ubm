package metric

import (
	"errors"
	"fmt"

	"../query"
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

func (msg *MetricGet) Receive(APIAddr string) (metrics map[string]int, err error) {
	resp, code := query.SendQuery(APIAddr, msg)
	if code != 200 {
		err = errors.New(fmt.Sprintf("response code: %d", code))
		return
	}

	var expectedMessage MetricGetResponse
	err = query.RestoreMessage(resp, "MetricGetResponse", &expectedMessage)
	metrics = expectedMessage.Metrics
	return
}
