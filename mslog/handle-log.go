package mslog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"com.gft.tsbo-training.src.go/common/device/implementation/devicemeasure"
	"com.gft.tsbo-training.src.go/common/ms-framework/microservice"
)

type LogEntry struct {
	microservice.Response
	devicemeasure.DeviceMeasure
}

var dbMutex sync.Mutex
var memoryDb map[string]*LogEntry = make(map[string]*LogEntry)

// ---------------------------------------------------------------------------
// httpGetLog ...

func (ms *MsLog) httpGetLog(w http.ResponseWriter, r *http.Request) (status int, contentLen int, msg string) {
	var response DBResponse
	dbMutex.Lock()
	msg = fmt.Sprintf("Got values from %d device(s).", len(memoryDb))
	InitDBResponse(&response, msg, &memoryDb, ms)
	dbMutex.Unlock()

	status = http.StatusOK
	ms.SetResponseHeaders("application/json; charset=utf-8", w, r)
	// w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// w.Header().Set("cid", ms.GetName())
	// w.Header().Set("version", ms.GetVersion())
	w.WriteHeader(status)
	contentLen = ms.Reply(w, response)
	return status, contentLen, msg
}

// ---------------------------------------------------------------------------
// httpPutLog ...

func (ms *MsLog) httpPutLog(w http.ResponseWriter, r *http.Request) (int, contentLen int, msg string) {
	status := http.StatusCreated
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		msg = fmt.Sprintf("Failed to read request body, error was '%s'!", err.Error())
		ms.SetResponseHeaders("", w, nil)
		w.WriteHeader(status)
		return http.StatusBadRequest, 0, msg
	}

	clientID := r.Header.Get("cid")

	if len(clientID) == 0 {
		clientID = r.Header.Get("name")
	}

	//	var measure devicemeasure.DeviceMeasure
	var logentry LogEntry
	var measure *devicemeasure.DeviceMeasure = &logentry.DeviceMeasure

	err = json.Unmarshal(body, &logentry)

	dbMutex.Lock()
	memoryDb[logentry.GetDeviceAddress()] = &logentry
	dbMutex.Unlock()

	if err != nil {
		msg = fmt.Sprintf("Failed to parse request, error was '%s'!", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return http.StatusInternalServerError, 0, msg
	}

	if ms.DBConnection == nil {
		status = http.StatusOK
		msg = fmt.Sprintf("Got value '%s' of '%s' at '%s' version '%s'.", measure.GetFormatted(), measure.GetDeviceAddress(), measure.GetStamp().Format("2006-01-02 15:04:05"), logentry.Version)
		var response LogResponse
		InitLogResponse(&response, msg, measure, ms)
		ms.SetResponseHeaders("application/json; charset=utf-8", w, r)
		// w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// w.Header().Set("cid", ms.GetName())
		// w.Header().Set("version", ms.GetVersion())
		w.WriteHeader(status)
		contentLen = ms.Reply(w, response)
		return status, contentLen, msg
	}

	if !ms.DBConnection.IsOpen() {
		err = ms.DBConnection.Open()
		if err != nil {
			msg = fmt.Sprintf("Failed to open database, error was '%s'!", err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return http.StatusInternalServerError, 0, msg
		}
		ms.GetLogger().Println("Opened database connection.")
	}

	err = ms.DBConnection.AddDevice(measure)

	if err != nil {
		msg = fmt.Sprintf("Failed to add device '%s', error was '%s'!", measure.GetDeviceAddress(), err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return http.StatusInternalServerError, 0, msg
	}

	err = ms.DBConnection.Update(measure)

	if err != nil {
		msg = fmt.Sprintf("Failed to update device '%s', error was '%s'!", measure.GetDeviceAddress(), err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return http.StatusInternalServerError, 0, msg
	}

	msg = fmt.Sprintf("Updated value '%s' of '%s' at '%s' in database.", measure.GetFormatted(), measure.GetDeviceAddress(), measure.GetStamp().Format("2006-01-02 15:04:05"))
	// response := NewLogResponse(msg, measure, ms)
	var response LogResponse
	InitLogResponse(&response, msg, measure, ms)
	ms.SetResponseHeaders("application/json; charset=utf-8", w, r)
	// w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// w.Header().Set("cid", ms.GetName())
	// w.Header().Set("version", ms.GetVersion())
	w.WriteHeader(status)
	contentLen = ms.Reply(w, response)
	return status, contentLen, msg

}
