package main

import (
	"github.com/Lamzin/go-ubm/server"
	"github.com/AlexeySpiridonov/goapp-config"
)

func main() {
	 //empty
}

func main_example_server() {

	serverAddress := config.Get("serverAddress")

	mongoDbHost := config.Get("dbHost")
	mongoDbName := config.Get("dbName")

	server := bserver.NewServer(serverAddress, mongoDbHost, mongoDbName)
	server.Run()

}
