package monitors

import (
	"wingoEDR/chainsaw"

	"github.com/Jeffail/gabs"
	"go.uber.org/zap"
)

func FullEventCheck() {
	_ = gabs.New()

	events, err := chainsaw.ScanAll()
	if err != nil {
		zap.S().Error("Chainsaw events were not scanned: ", err.Error())
	}

	for _, e := range events {
		chainsaw.RunEventResponse(e.ID)
	}

}

func RangedEventCheck(fromTimestamp string, toTimestamp string) {
	_ = gabs.New()

	events, err := chainsaw.ScanTimeRange(fromTimestamp, toTimestamp)
	if err != nil {
		zap.S().Error("Chainsaw events were not scanned: ", err.Error())
	}

	for _, e := range events {
		chainsaw.RunEventResponse(e.ID)
	}

}
