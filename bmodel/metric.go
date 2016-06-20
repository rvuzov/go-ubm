package bmodel

import (
        "gopkg.in/mgo.v2/bson"
)

func Inc(userId string, key string, value int) {
        _, err := Metrics.Upsert(
                bson.M{"id": userId},
                bson.M{"$inc": bson.M{key: value}},
        )
        refresh("bmodel", err)
}

func Log(userId string, key string, value interface{}) {
        limit := 128
        _, err := Metrics.Upsert(
                bson.M{"id": userId},
                bson.M{"$push": bson.M{"logs."+key: bson.M{"$each": []interface{}{value}, "$slice": -limit}}},
        )
        refresh("bmodel", err)
}
