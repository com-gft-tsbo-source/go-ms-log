package main

import (
	"os"
	"time"

	"com.gft.tsbo-training.src.go/ms-log/mslog"
)

// ###########################################################################
// ###########################################################################
// MAIN
// ###########################################################################
// ###########################################################################

// var _version string = "AAAA"

func main() {

	var ms mslog.MsLog
	mslog.InitFromArgs(&ms, os.Args, nil)

	go func() {

		if ms.DBConnection == nil {
			return
		}

		for ever := true; ever; ever = true {
			var err error

			if ms.DBConnection.IsOpen() {
				goto sleep
			}
			ms.GetLogger().Println("Trying to connect to database.")
			err = ms.DBConnection.Open()

			if err != nil {
				ms.GetLogger().Printf("Error: Failed to connect to database. Error was '%s'.\n", err.Error())
				goto sleep
			}

			ms.GetLogger().Println("Successfully connected to database.")

		sleep:
			time.Sleep(0 * time.Millisecond)
		}
	}()

	ms.Run()

}
