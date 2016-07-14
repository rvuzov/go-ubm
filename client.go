package main

import (
	"log"

	"github.com/Lamzin/go-ubm/client"
)

func main_example_client() {

	example_MetricPush()
	example_MetricGet()
	example_MetricFindUser()

	example_LogPush()

}

func example_LogPush() {
	client := bclient.NewClient("0.0.0.0:3001")

	lonlat := struct {
		Lon       string `json:"lon"`
		Lat       string `json:"lat"`
		TimeStamp int    `json:"timestamp"`
	}{"45.0321", "39.1654", 54645489}

	err := client.LogPush("good_user", "signin.geo", lonlat)
	if err != nil {
		log.Print(err.Error())
	} else {
		log.Print("Done!")
	}

}

func example_MetricPush() {
	client := bclient.NewClient("0.0.0.0:3001")

	err := client.MetricPush("good_user", "chat.text", 15)
	if err != nil {
		log.Print(err.Error())
	} else {
		log.Print("Done!")
	}

}

func example_MetricGet() {

	client := bclient.NewClient("0.0.0.0:3001")

	metrics, err := client.MetricGet("good_user", []string{"chat.text", "chat.doc", "geo.kiev"})
	if err != nil {
		log.Print(err.Error())
	} else {
		log.Print("Done!")
		log.Printf("%v\n", metrics)
	}

}

func example_MetricFindUser() {

	client := bclient.NewClient("0.0.0.0:3001")

	var cmps []bclient.Cmp
	cmps = append(
		cmps,
		bclient.Cmp{
			Metric:   "chat.text",
			Operator: ">",
			Value:    10,
		},
	)

	cmps = append(
		cmps,
		bclient.Cmp{
			Metric:   "chat.text",
			Operator: "<",
			Value:    20,
		},
	)

	result, err := client.MetricFindUsers(cmps)

	if err != nil {
		log.Print(err.Error())
	} else {
		log.Print("Done!")
		log.Printf("%v\n", result)
	}

}
