package main

import (
	"fmt"
	"wingoEDR/frontend"
	"wingoEDR/processes"
)

func main() {

	for {
		go frontend.QuickInterface()

		select {}
	}

	thing := processes.GetAllProcesses()

	thingNetworkConnections := thing[0].NetworkConnections

	fmt.Printf("%+v", thingNetworkConnections[0].NetType)

	//time.Sleep(1 * time.Minute)

}
