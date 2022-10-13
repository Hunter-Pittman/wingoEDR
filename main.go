package main

import (
	"time"
	"wingoEDR/frontend"
	"wingoEDR/logger"
	"wingoEDR/processes"
)

func main() {
	logger.InitLogger()
	for {
		//Interface
		go frontend.QuickInterface()
		// Process Analysis
		go getHashStatus("7a2278a9a74f49852a5d75c745ae56b80d5b4c16f3f6a7fdfd48cb4e2431c688", "sha256")
		go processes.GetAllProcesses()

		// Heartbeat

		time.Sleep(1 * time.Minute)
		select {}
	}
}
