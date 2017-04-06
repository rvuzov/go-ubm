package ubm

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

const (
	logsPushWorkersCount = 4
	logsChanSize         = 1000000
	logSizelimit         = 32
)

type (
	logs struct {
		lock  chan struct{}
		Queue chan string
		logs  map[string]*[]Log

		GetCalls          int64
		PushCalls         int64
		MongoUpsertCalls  int64
		pushLogsFrequency map[int]int64
	}

	Log struct {
		Key   string
		Value interface{}
	}
)

var Logs logs

func (l *logs) Init() {
	l.lock = make(chan struct{}, 1)
	l.Queue = make(chan string, logsChanSize)
	l.logs = make(map[string]*[]Log, 0)
	l.pushLogsFrequency = make(map[int]int64, 0)
	for i := 0; i < logsPushWorkersCount; i++ {
		go l.push()
	}
}

func (l *logs) GetPushMetricsFrequency() map[int]int64 {
	l.lock <- struct{}{}
	f := make(map[int]int64, len(l.pushLogsFrequency))
	for k, v := range l.pushLogsFrequency {
		f[k] = v
	}
	<-l.lock
	return f
}

func (l *logs) Push(userID string, key string, value interface{}) {
	l.lock <- struct{}{}
	if arr, ok := l.logs[userID]; ok {
		*arr = append(*arr, Log{Key: key, Value: value})
	} else {
		newArr := make([]Log, 1)
		newArr[0] = Log{Key: key, Value: value}
		l.logs[userID] = &newArr
		l.Queue <- userID
	}
	l.PushCalls++
	<-l.lock
}

func (l *logs) push() {
	for userID := range l.Queue {
		l.lock <- struct{}{}
		arr, ok := l.logs[userID]
		delete(l.logs, userID)
		l.MongoUpsertCalls++
		l.pushLogsFrequency[len(*arr)]++
		<-l.lock

		if !ok {
			loger.Errorf("user(%s) can't find metrics in map", userID)
			continue
		}

		unique := make(map[string]*[]interface{}, 0)
		for _, log := range *arr {
			if arrL, ok := unique[log.Key]; ok {
				*arrL = append(*arrL, log.Value)
			} else {
				newArrL := make([]interface{}, 1)
				newArrL[0] = log.Value
				unique[log.Key] = &newArrL
			}
		}

		push := make(map[string]interface{}, 0)
		for k, v := range unique {
			push["logs."+k] = bson.M{
				"$each":  *v,
				"$slice": -logSizelimit,
			}
		}

		_, err := Models.Upsert(
			bson.M{"id": userID},
			bson.M{
				"$push": push,
			})
		refresh("ubm", err)
	}
}

func (l *logs) Get(userID string, keys ...string) (answer map[string][]LogItem, err error) {
	var result map[string]interface{}
	answer = make(map[string][]LogItem, 0)

	project := bson.M{}
	for i, key := range keys {
		project[strconv.Itoa(i)] = "$logs." + key
	}

	err = Models.Pipe([]bson.M{
		bson.M{"$match": bson.M{"id": userID}},
		bson.M{"$project": project},
	}).One(&result)

	for i, key := range keys {
		if value, ok := result[strconv.Itoa(i)]; ok {
			if arr, ok := value.([]interface{}); ok {
				logArray := make([]LogItem, len(arr))
				for i, item := range arr {
					logArray[i] = LogItem{Item: item}
				}
				answer[key] = logArray
			}
		}
	}

	refresh("umb", err)
	l.GetCalls++
	return
}

type LogItem struct{
	Item interface{}
}

func (l LogItem) MustGetInt(key string) (i int) {
	if value, ok := l.get(key); ok {
		i, _ = value.(int)
	}
	return
}

func (l LogItem) MustGetString(key string) (s string) {
	if value, ok := l.get(key); ok {
		s, _ = value.(string)
	}
	return
}

func (l LogItem) get(key string) (interface{}, bool) {
	if m, ok := l.Item.(map[string]interface{}); ok {
		v, ok := m[key]
		return v, ok
	}
	return nil, false
}
