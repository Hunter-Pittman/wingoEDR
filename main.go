package main

import (
	"fmt"
	"wingoEDR/processes"
)

func main() {
	thing := processes.GetAllProcesses()

	thingNetworkConnections := thing[0].NetworkConnections

	fmt.Printf("%+v", thingNetworkConnections[0].NetType)
}
