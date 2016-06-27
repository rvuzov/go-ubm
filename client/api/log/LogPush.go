package log

import (
	"errors"
	"fmt"

	"../query"
	"../service"
)

type (
	LogPush struct {
		UserID string      `json:"userID"`
		Key    string      `json:"key"`
		Value  interface{} `json:"value"`
	}
)

func (msg *LogPush) Receive(APIAddr string) (err error) {
	resp, code := query.SendQuery(APIAddr, msg)
	if code != 200 {
		return errors.New(fmt.Sprintf("response code: %d", code))
	}

	var expectedMessage service.Success
	err = query.RestoreMessage(resp, "Success", &expectedMessage)
	return
}
