package metric

import (
	"errors"
	"fmt"

	"../query"
)

type (
	Cmp struct {
		Metric   string `json:"metric"`
		Operator string `json:"operator"`
		Value    int    `json:"value"`
	}
	MetricFindUsers []Cmp

	UserMetric struct {
		UserID  string         `json:"userID"`
		Metrics map[string]int `json:"metrics"`
	}
	MetricFindUsersResponse []UserMetric
)

func (msg *MetricFindUsers) Receive(APIAddr string) (answer MetricFindUsersResponse, err error) {
	resp, code := query.SendQuery(APIAddr, msg)
	if code != 200 {
		err = errors.New(fmt.Sprintf("response code: %d", code))
		return
	}

	var expectedMessage MetricFindUsersResponse
	err = query.RestoreMessage(resp, "MetricFindUsersResponse", &expectedMessage)
	answer = expectedMessage
	return
}
