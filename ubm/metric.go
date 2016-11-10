package ubm

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

const (
	metricsPushWorkersCount = 8
	metricsChanSize         = 1000000
)

type (
	metrics struct {
		lock    chan struct{}
		Queue   chan string
		Metrics map[string]*[]Metric

		GetCalls             int64
		PushCalls            int64
		MongoUpsertCalls     int64
		PushMetricsFrequency map[int]int64
	}

	Metric struct {
		Key   string
		Value int
	}
)

var Metrics metrics

func (m *metrics) Init() {
	m.lock = make(chan struct{}, 1)
	m.Queue = make(chan string, metricsChanSize)
	m.Metrics = make(map[string]*[]Metric, 0)
	m.PushMetricsFrequency = make(map[int]int64, 0)
	for i := 0; i < metricsPushWorkersCount; i++ {
		go m.push()
	}
}

func (m *metrics) Get(userID string, keys []string) (answer map[string]int, err error) {
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

	refresh("umb", err)
	m.GetCalls++
	return
}

func (m *metrics) Push(userID string, key string, value int) {
	m.lock <- struct{}{}
	if arr, ok := m.Metrics[userID]; ok {
		*arr = append(*arr, Metric{Key: key, Value: value})
	} else {
		newArr := make([]Metric, 1)
		newArr[0] = Metric{Key: key, Value: value}
		m.Metrics[userID] = &newArr
		m.Queue <- userID
	}
	<-m.lock
	m.PushCalls++
}

func (m *metrics) push() {
	for userID := range m.Queue {
		m.lock <- struct{}{}
		arr, ok := m.Metrics[userID]
		delete(m.Metrics, userID)
		<-m.lock

		if !ok {
			loger.Errorf("user(%s) can't find metrics in map", userID)
			continue
		}

		unique := make(map[string]int, 0)
		for _, metric := range *arr {
			unique[metric.Key] += metric.Value
		}

		_, err := Models.Upsert(
			bson.M{"id": userID},
			bson.M{"$inc": unique},
		)
		refresh("ubm", err)
		m.MongoUpsertCalls++
		m.PushMetricsFrequency[len(*arr)]++
	}
}
