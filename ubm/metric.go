package ubm

import (
	"strconv"

	"gopkg.in/mgo.v2/bson"
	"sync"
)

const (
	metricsPushWorkersCount = 8
	metricsChanSize         = 1000000
)

type (
	metrics struct {
		mutex   *sync.Mutex
		Queue   chan string
		Metrics map[string]*[]Metric
	}

	Metric struct {
		Key   string
		Value int
	}
)

var Metrics metrics

func (m *metrics) Init() {
	m.Queue = make(chan string, metricsChanSize)
	m.Metrics = make(map[string]*[]Metric, 0)
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
	return
}

func (m *metrics) Push(userID string, key string, value int) {
	m.mutex.Lock()
	if arr, ok := m.Metrics[userID]; ok {
		*arr = append(*arr, Metric{Key: key, Value: value})
	} else {
		newArr := make([]Metric, 1)
		newArr[0] = Metric{Key: key, Value: value}
		m.Metrics[userID] = &newArr
		m.Queue <- userID
	}
	m.mutex.Unlock()
}

func (m *metrics) push() {
	for userID := range m.Queue {
		m.mutex.Lock()
		arr, ok := m.Metrics[userID]
		delete(m.Metrics, userID)
		m.mutex.Unlock()

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
	}
}
