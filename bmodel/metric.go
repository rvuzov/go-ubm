package bmodel

import (
	"gopkg.in/mgo.v2/bson"
)

type (
	UMetric struct {
		UserID string `json:"user"`
		Key    string `json:"key"`
		Value  int    `json:"value"`
	}
)

func (m *UMetric) Push() (err error) {
	_, err = Metrics.Upsert(
		bson.M{"id": (*m).UserID},
		bson.M{"$inc": bson.M{(*m).Key: (*m).Value}},
	)
	refresh("bmodel", err)
	return
}
