package bmodel

import (
	"gopkg.in/mgo.v2/bson"
)

type (
	metrics struct{}

	UMetric struct {
		UserID string `json:"user"`
		Key    string `json:"key"`
		Value  int    `json:"value"`
	}
)

var (
	Metrics metrics
)

func (_ metrics) Push(m *UMetric) (err error) {
	_, err = Models.Upsert(
		bson.M{"id": (*m).UserID},
		bson.M{"$inc": bson.M{(*m).Key: (*m).Value}},
	)
	refresh("bmodel", err)
	return
}
