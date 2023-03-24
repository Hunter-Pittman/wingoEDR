package monitors

import (
	"wingoEDR/chainsaw"

	"go.uber.org/zap"
)

func FullEventReportCheck() {

	events, err := chainsaw.ScanAll()
	if err != nil {
		zap.S().Error("Chainsaw events were not scanned: ", err.Error())
	}

	for _, e := range events {
		chainsaw.RunReportEventResponse(e)
	}

}

func RangedEventReportCheck(fromTimestamp string, toTimestamp string) {

	events, err := chainsaw.ScanTimeRange(fromTimestamp, toTimestamp)
	if err != nil {
		zap.S().Error("Chainsaw events were not scanned: ", err.Error())
	}

	for _, e := range events {
		chainsaw.RunReportEventResponse(e)
	}

}

func FullEventActionCheck() {

	events, err := chainsaw.ScanAll()
	if err != nil {
		zap.S().Error("Chainsaw events were not scanned: ", err.Error())
	}

	for _, e := range events {
		chainsaw.RunActionEventResponse(e.ID)
	}

}

func RangedEventActionCheck(fromTimestamp string, toTimestamp string) {

	events, err := chainsaw.ScanTimeRange(fromTimestamp, toTimestamp)
	if err != nil {
		zap.S().Error("Chainsaw events were not scanned: ", err.Error())
	}

	for _, e := range events {
		chainsaw.RunActionEventResponse(e.ID)
	}

}
