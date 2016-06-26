package main

import (
	"./server"
)

func main() {

	server := bserver.New("0.0.0.0:3001")
	server.Run()

}
