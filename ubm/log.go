package ubm

import (
	"gopkg.in/mgo.v2/bson"
	"sync"
)

const (
	logsPushWorkersCount = 4
	logsChanSize         = 1000000
	logSizelimit         = 32
)

type (
	logs struct {
		mutex *sync.Mutex
		Queue chan string
		Logs  map[string]*[]Log
	}

	Log struct {
		Key   string
		Value interface{}
	}
)

var Logs logs

func (l *logs) Init() {
	l.mutex = &sync.Mutex{}
	l.Queue = make(chan string, logsChanSize)
	l.Logs = make(map[string]*[]Log, 0)
	for i := 0; i < logsPushWorkersCount; i++ {
		go l.push()
	}
}

func (l *logs) Push(userID string, key string, value interface{}) {
	l.mutex.Lock()
	if arr, ok := l.Logs[userID]; ok {
		*arr = append(*arr, Log{Key: key, Value: value})
	} else {
		newArr := make([]Log, 1)
		newArr[0] = Log{Key: key, Value: value}
		l.Logs[userID] = &newArr
		l.Queue <- userID
	}
	l.mutex.Unlock()
}

func (l *logs) push() {
	for userID := range l.Queue {
		l.mutex.Lock()
		arr, ok := l.Logs[userID]
		delete(l.Logs, userID)
		l.mutex.Unlock()

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
