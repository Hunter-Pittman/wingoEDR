package chainsaw

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"go.uber.org/zap"
)

func FullEventCheck() {
	_ = gabs.New()

	events, err := ScanAll()
	if err != nil {
		zap.S().Error("Chainsaw events were not scanned: ", err.Error())
	}

	fmt.Printf("%+v", events[0])

}
