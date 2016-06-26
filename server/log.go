package bserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	bmodel "../bmodel"
	"github.com/julienschmidt/httprouter"
)

type (
	ULog struct {
		UserID string      `json:"user"`
		Key    string      `json:"key"`
		Value  interface{} `json:"value"`
	}
)

func pushLog(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		writeResponse(w, 400, "Bad requests: "+err.Error())
		return
	}

	userID := r.Form.Get("userID")
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	if userID == "" || key == "" || value == "" {
		writeResponse(w, 400, "Bad requests: empty parameter")
		return
	}

	var valuei interface{}
	err := json.Unmarshal([]byte(value), &valuei)
	if err != nil {
		writeResponse(w, 400, fmt.Sprintf("Bad value parameter: %v", err.Error()))
		return
	}

	ulog := bmodel.ULog{
		UserID: userID,
		Key:    key,
		Value:  valuei,
	}
	err = ulog.Push()
	if err != nil {
		writeResponse(w, 500, fmt.Sprintf("Shit happens: %v", err.Error()))
		return
	}

	writeResponse(w, 200, "Done")
}
