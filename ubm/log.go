package ubm

import (
	"gopkg.in/mgo.v2/bson"
)

type (
	logs struct{}
)

var (
	limit = 32
	Logs  logs
)

func (_ logs) Push(userID string, key string, value interface{}) (err error) {
	_, err = Models.Upsert(
		bson.M{"id": userID},
		bson.M{
			"$push": bson.M{
				"logs." + key: bson.M{
					"$each":  []interface{}{value},
					"$slice": -limit,
				},
			},
		})
	refresh("ubm", err)
	return
}
