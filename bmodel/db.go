package bmodel

import (
        "github.com/AlexeySpiridonov/goapp-config"
	"github.com/op/go-logging"
	"gopkg.in/mgo.v2"
)

const (
	metricDbName = "user"
)

var (
	log = logging.MustGetLogger("bmodel")

	context Context

	Metrics *mgo.Collection
)

type Context struct {
	Session *mgo.Session
	Db      *mgo.Database
}

func Get() Context {
	return context
}

func Init() (*mgo.Session, error) {
	log.Info("Connect to DB: " + config.Get("dbHost") + " " + config.Get("dbName"))
	mongo, err := mgo.Dial(config.Get("dbHost"))
	if err != nil {
		log.Panic("Cant't connect to mongoDB. Server is stopped")
	}
	log.Info("DB ok")
	set(mongo, mongo.DB(config.Get("dbName")))
	return mongo, err
}

func set(session *mgo.Session, db *mgo.Database) {
	context = Context{session, db}

	Metrics = context.Db.C(metricDbName)
}

func refresh(source string, err error) {
	if err == nil {
		return
	}

	if err.Error() == "not found" {
		log.Notice(source + " " + err.Error())
	} else {
		log.Error(source + " " + err.Error())
	}

	if err.Error() == "EOF" {
		log.Warning("DB connect autoRefresh")
		context.Session.Refresh()
	}
}
