package bmodel

import (
	"gopkg.in/mgo.v2/bson"
)

type (
	ULog struct {
		UserID string      `json:"user"`
		Key    string      `json:"key"`
		Value  interface{} `json:"value"`
	}
)

var (
	limit = 128
)

func (l *ULog) Push() (err error) {
	_, err = Metrics.Upsert(
		bson.M{"id": (*l).UserID},
		bson.M{
			"$push": bson.M{
				"logs." + (*l).Key: bson.M{
					"$each":  []interface{}{(*l).Value},
					"$slice": -limit,
				},
			},
		})
	refresh("bmodel", err)
	return
}
