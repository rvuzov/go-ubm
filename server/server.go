package bserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	bmodel "../bmodel"
	"./api"
	"github.com/julienschmidt/httprouter"
)

type (
	BServer struct {
		Addr        string
		MongoDbHost string
		MongoDbName string
	}
)

func NewServer(addr string, mongoDbHost string, mongoDbName string) BServer {
	server := BServer{
		Addr:        addr,
		MongoDbHost: mongoDbHost,
		MongoDbName: mongoDbName,
	}
	return server
}

func (server *BServer) Run() {

	mongoSession, err := bmodel.Init((*server).MongoDbHost, (*server).MongoDbName)
	defer mongoSession.Close()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	router := httprouter.New()

	router.POST("/", apiController)

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

func apiController(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	query := r.FormValue("query")
	resp, code := api.Process(query)
	writeResponse(w, code, resp)

	log.Print(query)
}
