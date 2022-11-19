package main

import (
	"time"
	"wingoEDR/common"
	"wingoEDR/honeymonitor"
	"wingoEDR/logger"
)

func main() {
	logger.InitLogger()
	go heartbeatLoop()
	go inventoryLoop()
	go objectMonitoring()
}

func inventoryLoop() {
	ticker := time.NewTicker(20 * time.Second)

	for _ = range ticker.C {
		common.PostInventory()
	}
}

func heartbeatLoop() {
	ticker := time.NewTicker(1 * time.Minute)

	for _ = range ticker.C {
		common.PostInventory()
	}
}

func objectMonitoring() {
	ticker := time.NewTicker(10 * time.Second)

	for _ = range ticker.C {
		montiorPath := "C:\\Users\\hunte\\Documents\\honey_monitor_test"

		var pathVector []string

		pathVector = append(pathVector, montiorPath)

		honeymonitor.CreateDirMonitor(pathVector)
	}
}
