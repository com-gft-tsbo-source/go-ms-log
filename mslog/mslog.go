package mslog

import (
	"flag"
	"time"

	"com.gft.tsbo-training.src.go/common/device/devicedb"
	"com.gft.tsbo-training.src.go/common/ms-framework/microservice"
)

// ###########################################################################
// ###########################################################################
// MsLog
// ###########################################################################
// ###########################################################################

// MsLog ...
type MsLog struct {
	microservice.MicroService
	starttime    time.Time
	DBConnection devicedb.IConnection
}

// ###########################################################################

// InitMsLogFromArgs ...
func InitFromArgs(ms *MsLog, args []string, flagset *flag.FlagSet) *MsLog {

	var db devicedb.IConnection

	if flagset == nil {
		flagset = flag.NewFlagSet("ms-log", flag.PanicOnError)
	}

	microservice.InitFromArgs(&ms.MicroService, args, flagset, ms.DefaultHandler())

	if len(ms.GetDBName()) > 0 {
		db = devicedb.NewDatabase(ms.GetDBName(), "measurements", ms.GetName())
		_, isNil := db.(*devicedb.NilConnection)

		if !isNil {
			ms.GetLogger().Printf("Got database configuration '%s'.\n", ms.GetDBName())
		} else {
			ms.GetLogger().Printf("Bad database configuration '%s', ignoring it.\n", ms.GetDBName())
			db = nil
		}
	} else {
		ms.GetLogger().Println("No database configured.")
	}

	ms.DBConnection = db
	ms.starttime = time.Now()

	if ms.DBConnection != nil {
		ms.DBConnection.Open()
		ms.GetLogger().Printf("Opened database '%s'.\n", ms.GetDBName())
	}

	handlerLog := ms.DefaultHandler()
	handlerLog.Put = ms.httpPutLog
	handlerLog.Get = ms.httpGetLog
	ms.AddHandler("/log", handlerLog)
	return ms
}
