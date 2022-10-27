package main

import (
	"time"
	"wingoEDR/logger"
)

func main() {
	logger.InitLogger()
	for {
		//Interface
		//go frontend.QuickInterface()

		//Inventory()

		// Yara Scan
		//yara.DirYaraScan()
		// Get Process Analysis
		//go processes.GetAllProcesses()

		// Kapersky API Hash check
		//go GetHashStatus("7a2278a9a74f49852a5d75c745ae56b80d5b4c16f3f6a7fdfd48cb4e2431c688", "sha256")

		// Heartbeat

		time.Sleep(1 * time.Minute)
		select {}
	}
}
