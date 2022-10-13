package main

import (
	"fmt"
	"wingoEDR/logger"
	"wingoEDR/processes"
)

func main() {
	logger.InitLogger()

	/*
		for {
			go frontend.QuickInterface()

			thing := processes.GetAllProcesses()

			thingNetworkConnections := thing[0].NetworkConnections

			fmt.Printf("%+v", thingNetworkConnections[0].NetType)

			select {}
		}
	*/
	//time.Sleep(1 * time.Minute)

	thing := processes.GetAllProcesses()
	fmt.Printf("%+v", thing)
}
