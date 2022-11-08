package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

func main() {
	//logger.InitLogger()
	logger, _ := zap.NewProduction()
	logger.Warn("log message test")
	for {
		//Interface
		//go frontend.QuickInterface()

		//Inventory
		thing := GetInventory()

		fmt.Printf("%v", thing)

		// Yara Scan
		// thing, err := yara.YaraScan("C:\\Users\\hunte\\Documents\\repos\\wingoEDR\\yara_rules\\fileID\\", "C:\\Users\\hunte\\Pictures")
		// if err != nil {
		// 	zap.S().Error(err)
		// }

		// for i := range thing {
		// 	fmt.Println(i)
		// 	fmt.Printf("%v", thing[0].Rule)
		// }

		// Get Process Analysis
		//go processes.GetAllProcesses()

		// Kapersky API Hash check
		//go GetHashStatus("7a2278a9a74f49852a5d75c745ae56b80d5b4c16f3f6a7fdfd48cb4e2431c688", "sha256")

		// Heartbeat

		time.Sleep(1 * time.Minute)
		select {}
	}
}
