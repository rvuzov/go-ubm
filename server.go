package main

import (
	"./server"
	"github.com/AlexeySpiridonov/goapp-config"
)

func main() {

	serverAddress := config.Get("serverAddress")

	mongoDbHost := config.Get("dbHost")
	mongoDbName := config.Get("dbName")

	server := bserver.New(serverAddress, mongoDbHost, mongoDbName)
	server.Run()

}
