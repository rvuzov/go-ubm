package bserver

import (
	"fmt"
	"net/http"
	"strconv"

	bmodel "../bmodel"
	"github.com/julienschmidt/httprouter"
)

type (
	UMetric struct {
		UserID string `json:"user"`
		Key    string `json:"key"`
		Value  int    `json:"value"`
	}
)

func pushMetric(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	valuei, err := strconv.Atoi(value)
	if err != nil {
		writeResponse(w, 400, fmt.Sprintf("Bad requests: %v", err.Error()))
		return
	}

	umetric := bmodel.UMetric{
		UserID: userID,
		Key:    key,
		Value:  valuei,
	}
	err = bmodel.Metrics.Push(&umetric)
	if err != nil {
		writeResponse(w, 500, fmt.Sprintf("Shit happens: %v", err.Error()))
		return
	}

	writeResponse(w, 200, "Done")
}
