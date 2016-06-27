package bmodel

import "gopkg.in/mgo.v2/bson"

type (
	metrics struct{}

	CompareStatement struct {
		Key      string `json:"key"`
		Operator string `json:"operator"`
		Value    int    `json:"value"`
	}
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

func (_ metrics) GetKeys(userID string, keys []string) (result []interface{}, err error) {
	project := bson.M{
		"id": "$id",
	}
	for _, key := range keys {
		project[key] = "$" + key
	}

	err = Models.Pipe([]bson.M{
		bson.M{"$match": bson.M{"id": userID}},
		bson.M{"$project": project},
	}).All(&result)

	refresh("bmodel", err)
	return
}
