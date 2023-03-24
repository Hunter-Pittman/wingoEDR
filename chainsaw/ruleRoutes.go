package chainsaw

import (
	"encoding/json"
	"os"
	"time"
	"wingoEDR/apis/serialscripter"
	"wingoEDR/usermanagement"

	"go.uber.org/zap"
)

// this function reports the event to serialscripter
func RunReportEventResponse(event ChainsawEvent) {
	newChainsawIncident(event)
}

func newChainsawIncident(event ChainsawEvent) {
	zap.S().Infof("New chainsaw incident detected: RuleName: %v Time: %v,", event.RuleName, event.Timestamp)
	payload := event

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		zap.S().Error("Error marshalling payload: ", err)
	}

	newIncident := serialscripter.Incident{
		Name:        event.RuleName,
		CurrentTime: time.Now().String(),
		User:        usermanagement.GetLastLoggenOnUser(),
		Severity:    event.Level,
		Payload:     string(jsonPayload),
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	newAlert := serialscripter.Alert{
		Host:     hostname,
		Incident: newIncident,
	}
	serialscripter.IncidentAlert(newAlert)
}

var functionMap = map[string]func(){
	"78d5cab4-557e-454f-9fb9-a222bd0d5edc": Ident78d5cab4557e454f9fb9a222bd0d5edc,
}

var count int

func RunActionEventResponse(id string) {
	function, exists := functionMap[id]
	if exists {
		function()
	} else {
		//fmt.Println("Function does not exist.")
		return
	}
}

// RDP Connection
func Ident78d5cab4557e454f9fb9a222bd0d5edc() {
	println("Event response for 78d5cab4-557e-454f-9fb9-a222bd0d5edc")
}
