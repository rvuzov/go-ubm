package metric

import (
	"errors"
	"fmt"

	"../query"
	"../service"
)

type (
	MetricPush struct {
		UserID string `json:"userID"`
		Key    string `json:"key"`
		Value  int    `json:"value"`
	}
)

func (msg *MetricPush) Receive(APIAddr string) (err error) {
	resp, code := query.SendQuery(APIAddr, msg)
	if code != 200 {
		return errors.New(fmt.Sprintf("response code: %d", code))
	}
	fmt.Println(resp)

	var expectedMessage service.Success
	err = query.RestoreMessage(resp, "Success", &expectedMessage)
	return
}
