package mslog

import (
	"github.com/com-gft-tsbo-source/go-common/device/implementation/devicedescriptor"
	"github.com/com-gft-tsbo-source/go-common/device/implementation/devicemeasure"
	"github.com/com-gft-tsbo-source/go-common/ms-framework/microservice"
)

// ###########################################################################
// ###########################################################################
// MsLog Response - Put
// ###########################################################################
// ###########################################################################

// LogResponse Encapsulates the reploy of ms-measure
type LogResponse struct {
	microservice.Response
	devicedescriptor.DeviceDescriptor
}

// ###########################################################################

// InitLogResponse Constructor of a response of ms-measure
func InitLogResponse(r *LogResponse, status string, obj devicemeasure.IDeviceMeasure, ms *MsLog) {
	microservice.InitResponseFromMicroService(&r.Response, ms, 200, status)
	devicedescriptor.InitFromDeviceDescriptor(&r.DeviceDescriptor, obj)
}

// ###########################################################################
// ###########################################################################
// MsLog DBResponse - Get
// ###########################################################################
// ###########################################################################

// DBResponse ...
type DBResponse struct {
	microservice.Response
	Entries map[string]LogEntry `json:entries`
}

// InitDBResponse ...
func InitDBResponse(r *DBResponse, status string, entries *map[string]*LogEntry, ms *MsLog) {
	microservice.InitResponseFromMicroService(&r.Response, ms, 200, status)
	r.Entries = make(map[string]LogEntry)

	for key, value := range *entries {
		r.Entries[key] = *value
	}

}
