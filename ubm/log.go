package ubm

import "gopkg.in/mgo.v2/bson"

const (
	logsPushWorkersCount = 4
	logsChanSize         = 1000000
	logSizelimit         = 32
)

type (
	logs struct {
		lock              chan struct{}
		Queue             chan string
		logs              map[string]*[]Log

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
