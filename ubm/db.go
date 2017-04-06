package ubm

import (
	"strings"

	"github.com/op/go-logging"
	"gopkg.in/mgo.v2"
)

const (
	modelDbName = "user"
)

type (
	Context struct {
		Session *mgo.Session
		Db      *mgo.Database
	}
)

var (
	loger   = logging.MustGetLogger("ubm")
	context Context
	Models  *mgo.Collection
)

func Init(dbHost string, dbName string) (*mgo.Session, error) {
	loger.Infof("Connect to DB: %s %s", dbHost, dbName)
	mongoSession, err := mgo.Dial(dbHost)
	if err != nil {
		loger.Panic("Cant't connect to mongoDB. Server is stopped")
	}
	loger.Info("DB ok")

	mongoSession.SetMode(mgo.Monotonic, true)
	mongoSession.SetSafe(nil)
	mongoSession.Fsync(false)

	context = Context{
		Session: mongoSession,
		Db:      mongoSession.DB(dbName),
	}

	Models = context.Db.C(modelDbName)

	Metrics.Init()
	loger.Info("init metrics: OK")
	Logs.Init()
	loger.Info("init logs: OK")

	return mongoSession, err
}

func refresh(source string, err error) {
	if err == nil {
		return
	}

	if err.Error() == "not found" {
		loger.Noticef("%s %s", source, err.Error())
	} else {
		loger.Errorf("%s %s", source, err.Error())
	}

	if strings.Contains(err.Error(), "timeout") {
		loger.Warning("UBM refresh by timeout", source, err)
		context.Session.Refresh()
	}

	if strings.Contains(err.Error(), "connection reset by peer") {
		loger.Warning("UBM refresh by reset by peer", source, err)
		context.Session.Refresh()
	}

	if err.Error() == "EOF" {
		loger.Warning("UBM refresh by EOF", source, err)
		context.Session.Refresh()
	}
}
