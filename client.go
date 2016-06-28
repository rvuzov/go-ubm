package main

import (
	"log"

	"./client"
)

func main() {

	client := bclient.NewClient("0.0.0.0:3001")

	// metrics, err := client.MetricGet("good user", []string{"chat.text", "chat.doc", "geo.kiev"})
	// if err != nil {
	// 	log.Print(err.Error())
	// } else {
	// 	log.Print("Done!")
	// 	log.Printf("%v\n", metrics)
	// }

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
