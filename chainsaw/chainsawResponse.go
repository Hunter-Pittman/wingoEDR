package chainsaw

import (
	"fmt"

	"github.com/Jeffail/gabs"
)

func FullEventCheck() {
	_ = gabs.New()

	events, err := ScanAll()
	if err != nil {
		//log.Error("Error scanning all events: ", err.Error())
	}

	childern, _ := events.Children()

	for _, child := range childern {
		fmt.Println(child.Data().(string))
	}
}
