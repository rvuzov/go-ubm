package bmodel

import (
	"gopkg.in/mgo.v2/bson"
)

type (
	metrics struct{}
)

var (
	Metrics metrics
)

func (_ metrics) Push(userID string, key string, value int) (err error) {
	_, err = Models.Upsert(
		bson.M{"id": userID},
		bson.M{"$inc": bson.M{key: value}},
	)
	refresh("bmodel", err)
	return
}
