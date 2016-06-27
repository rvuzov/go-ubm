package main

import (
	"log"

	"./client"
)

func main() {

	client := bclient.NewClient("0.0.0.0:3001")

	metrics, err := client.GetMetric("good user", []string{"chat.text", "chat.doc", "geo.kiev"})
	if err != nil {
		log.Print(err.Error())
	} else {
		log.Print("Done!")
		log.Printf("%v\n", metrics)
	}

}
