package main

import (
	"sync"
	"time"
	"wingoEDR/common"
	"wingoEDR/honeymonitor"
	"wingoEDR/logger"
)

func main() {
	logger.InitLogger()

	var wg sync.WaitGroup
	wg.Add(3)

	go heartbeatLoop()
	go inventoryLoop()
	go objectMonitoring()

	wg.Wait()
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
		common.HeartBeat()
	}
}

func objectMonitoring() {
	ticker := time.NewTicker(10 * time.Second)

	for _ = range ticker.C {
		honeymonitor.CreateDirMonitor(common.GetHoneyPaths())
	}
}
