package main

import (
	"time"
	"wingoEDR/common"
	"wingoEDR/honeydirectory"
	"wingoEDR/honeytoken"
	"wingoEDR/logger"
)

func main() {
	logger.InitLogger()
	for {

		//go inventoryLoop()
		//go honeytokenLoop()
		//go honeydirectoryLoop()

		//processes.WindowsFindDeviousCmdParams2()

		time.Sleep(1 * time.Minute)
		select {}
	}
}

func inventoryLoop() {
	ticker := time.NewTicker(10 * time.Second)

	for _ = range ticker.C {
		common.PostInventory()
	}
}

func honeytokenLoop() {
	ticker := time.NewTicker(10 * time.Second)

	for _ = range ticker.C {
		honeytoken.MonitorHoneyFile("C:\\Windows\\setupact.log")
	}
}

func honeydirectoryLoop() {
	ticker := time.NewTicker(10 * time.Second)

	for _ = range ticker.C {
		honeydirectory.MonitorHoneyDirectory("C:\\Windows\\setupact.log")
	}
}
