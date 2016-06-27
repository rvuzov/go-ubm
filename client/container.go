package bclient

import (
	"encoding/json"
	"reflect"
)

type (
	APIContainer struct {
		Type    string `json:"type"`
		Message string `json:"message`
	}
)

func NewAPIContainer(msg interface{}) (container APIContainer, err error) {
	typeName := reflect.TypeOf(msg).Name()
	if typeName == "" {
		typeName = reflect.TypeOf(msg).Elem().Name()
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return
	}

	container = APIContainer{
		Type:    typeName,
		Message: string(bytes[:]),
	}
	return
}
