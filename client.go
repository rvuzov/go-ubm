package main

import (
	"log"

	"./client"
)

func main() {

	client := bclient.NewClient("0.0.0.0:3001")

	err := client.PushMetric("good user", "chat.text", 45)
	if err != nil {
		log.Print(err.Error())
	} else {
		log.Print("Done!")
	}

}
