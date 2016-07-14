package ubm

import (
	"errors"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

type (
	metrics struct{}

	UserMetric struct {
		UserID  string         `json:"userID"`
		Metrics map[string]int `json:"metrics"`
	}

	CompareStatement struct {
		Metric   string `json:"metric"`
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

func (_ metrics) Get(userID string, keys []string) (answer map[string]int, err error) {
	var result map[string]int
	answer = make(map[string]int)

	project := bson.M{}
	for i, key := range keys {
		project[strconv.Itoa(i)] = "$" + key // ugly hack
	}

	err = Models.Pipe([]bson.M{
		bson.M{"$match": bson.M{"id": userID}},
		bson.M{"$project": project},
	}).One(&result)

	for i, key := range keys {
		if value, ok := result[strconv.Itoa(i)]; ok {
			answer[key] = value
		}
	}

	refresh("bmodel", err)
	return
}

func (_ metrics) FindUsers(statements []CompareStatement) (answer []UserMetric, err error) {

	var mongoStatements []bson.M
	for _, statement := range statements {
		mongoCmp, err := statement.ToMongoComparison()
		if err != nil {
			return answer, err
		}
		mongoStatements = append(mongoStatements, mongoCmp)
	}
	match := bson.M{"$and": mongoStatements}

	project := bson.M{
		"id": "$id",
	}
	for i, statement := range statements {
		project[strconv.Itoa(i)] = "$" + statement.Metric
	}

	var results []bson.M
	err = Models.Pipe(
		[]bson.M{
			bson.M{"$match": match},
			bson.M{"$project": project},
		},
	).All(&results)
	if err != nil {
		return
	}

	for _, result := range results {
		var tmpMetric UserMetric
		tmpMetric.Metrics = make(map[string]int)

		v, _ := result["id"]
		tmpMetric.UserID, _ = v.(string)

		for i, statement := range statements {
			v, _ := result[strconv.Itoa(i)]
			tmpMetric.Metrics[statement.Metric], _ = v.(int)
		}
		answer = append(answer, tmpMetric)
	}

	return
}

func (statement *CompareStatement) ToMongoComparison() (result bson.M, err error) {
	var operator string

	switch statement.Operator {
	case "=":
		operator = "$eq"
	case "!=":
		operator = "$ne"
	case ">":
		operator = "$gt"
	case "<":
		operator = "$lt"
	case ">=":
		operator = "$gte"
	case "<=":
		operator = "$lte"
	default:
		err = errors.New("unknown compare operator: " + statement.Operator)
		return
	}

	result = bson.M{
		(*statement).Metric: bson.M{
			operator: (*statement).Value,
		},
	}

	return
}
