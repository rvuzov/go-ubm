package bserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	bmodel "../bmodel"
	"github.com/julienschmidt/httprouter"
)

type (
	BServer struct {
		Addr string
	}
)

func New(addr string) BServer {
	server := BServer{
		Addr: addr,
	}
	return server
}

func (server *BServer) Run() {

	mongo, err := bmodel.Init()
	defer mongo.Close()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	router := httprouter.New()

	router.GET("/push.metric", pushMetric)
	router.GET("/push.log", pushLog)

	log.Printf("Run go-ubm server on http://%s", (*server).Addr)
	log.Fatal(http.ListenAndServe((*server).Addr, router))
}

func writeResponse(w http.ResponseWriter, code int, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.WriteHeader(code)

	jsonResult, _ := json.Marshal(result)
	fmt.Fprintf(w, "%s", jsonResult)
}
