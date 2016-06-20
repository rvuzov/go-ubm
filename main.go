package main

import (
	"./bmodel"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
        ready = make(chan struct{})
)

func main() {

	mongo, err := bmodel.Init()
	defer mongo.Close()
	if err != nil {
		return
	}

        n := 100
	for i := 0; i < n; i++ {
		go test()
	}

        for i := 0; i < n; i++ {
		<-ready
	}

}

func test() {
	for i := 0; i < 100; i++ {
                var id = RandStringBytes(20)
		bmodel.Inc(id, "chat.text", 25)
		bmodel.Inc(id, "chat.doc", 1)
		bmodel.Inc(id, "chat.image", 1)

		bmodel.Inc(id, "notification.push", 1)
		bmodel.Inc(id, "notification.popup", 2)
		bmodel.Inc(id, "notification.email", 15)

		bmodel.Log(id, "chat.text", "text message")
		bmodel.Log(id, "chat.doc", 15)
		bmodel.Log(id, "chat.image", map[string]interface{}{"url": "https://Logress.com/img15.jpg", "name": "img"})

		bmodel.Log(id, "notification.push", 1)
		bmodel.Log(id, "notification.popup", 2)
		bmodel.Log(id, "notification.email", map[string]interface{}{"subject": "something happens", "from": "bot"})
	}
        ready<-struct{}{}
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
