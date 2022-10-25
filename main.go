package main

import (
	"time"
	"wingoEDR/logger"
	"wingoEDR/yara"
)

func main() {
	logger.InitLogger()
	for {
		//Interface
		//go frontend.QuickInterface()
		yara.DirYaraScan()
		// Process Analysis

		// Heartbeat

		// yara scan

		time.Sleep(1 * time.Minute)
		select {}
	}
}
