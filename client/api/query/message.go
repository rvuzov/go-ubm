package query

import (
	"encoding/json"
	"errors"

	"github.com/Lamzin/go-ubm/client/api/service"
)

func RestoreMessage(s string, expected string, msg interface{}) (err error) {
	container, err := RestoreContainer(s)
	if err != nil {
		return
	}

	switch container.Type {
	case expected:
		err = json.Unmarshal([]byte(container.Message), &msg)
	case "Error":
		var APIerr service.Error
		err = json.Unmarshal([]byte(container.Message), &APIerr)
		if err != nil {
			return
		}
		err = errors.New(APIerr.Error)
	default:
		err = errors.New("unexpected type: " + container.Type)
	}
	return
}
