package chainsaw

import (
	"github.com/Jeffail/gabs"
	"go.uber.org/zap"
)

func FullEventCheck() {
	_ = gabs.New()

	events, err := ScanAll()
	if err != nil {
		zap.S().Error("Chainsaw events were not scanned: ", err.Error())
	}

	for _, e := range events {
		RunEventResponse(e.ID)
	}

}

func RangedEventCheck(fromTimestamp string, toTimestamp string) {
	_ = gabs.New()

	events, err := ScanTimeRange(fromTimestamp, toTimestamp)
	if err != nil {
		zap.S().Error("Chainsaw events were not scanned: ", err.Error())
	}

	for _, e := range events {
		RunEventResponse(e.ID)
	}

}
