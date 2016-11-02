package api

import (
	"encoding/json"
	"errors"

	"reflect"

	"github.com/Lamzin/go-ubm/server/api/log"
	"github.com/Lamzin/go-ubm/server/api/metric"
	"github.com/Lamzin/go-ubm/server/api/service"
)

type (
	APIContainer struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	}
)

func Process(query string) (string, int) {
	var container APIContainer
	err := json.Unmarshal([]byte(query), &container)
	if err != nil {
		return "error: " + err.Error(), 400
	}
	return container.Process()
}

func (c *APIContainer) Process() (string, int) {
	var msg APIMessage
	var resp interface{}

	switch (*c).Type {
	case "MetricPush":
		msg = new(metric.MetricPush)
	case "MetricGet":
		msg = new(metric.MetricGet)
	case "LogPush":
		msg = new(log.LogPush)
	default:
		err := errors.New("unknown API method")
		resp = service.NewError(err.Error())
		return responseToString(resp), 500
	}

	err := json.Unmarshal([]byte((*c).Message), &msg)
	if err != nil {
		resp = service.NewError(err.Error())
		return responseToString(resp), 500
	}

	resp = msg.Receive()
	return responseToString(resp), 200
}

func responseToString(msg interface{}) string {
	typeName := reflect.TypeOf(msg).Name()
	if typeName == "" {
		typeName = reflect.TypeOf(msg).Elem().Name()
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return "error: " + err.Error()
	}

	container := APIContainer{
		Type:    typeName,
		Message: string(bytes[:]),
	}

	bytes, err = json.Marshal(container)
	if err != nil {
		return "error: " + err.Error()
	}

	return string(bytes[:])
}
