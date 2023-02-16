package main

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"sync"
	"time"
	"wingoEDR/common"
	"wingoEDR/honeymonitor"
	"wingoEDR/logger"

	"github.com/fatih/color"
	"go.uber.org/zap"
)

func main() {
	// Command line args
	defaultConfigPath := "C:\\Users\\FORENSICS\\AppData\\Roaming\\wingoEDR\\config.json"
	configPtr := flag.String("config", defaultConfigPath, "Provide path to the config file")

	isStandalone := flag.Bool("standalone", true, "If serial scripter is not available then it outputs datga in local csv")

	flag.Parse()

	isWindowsPath := common.VerifyWindowsPath(*configPtr)

	if !isWindowsPath { // Errors out on a "C:\" path needs to be fixed
		color.Red("ERROR	The entered output is not a Windows path!")
		os.Exit(1)
	} else {
		if _, err := os.Stat(*configPtr); os.IsNotExist(err) {
			color.Red("ERROR	Windows path does not exist!")
			os.Exit(1)
		}
		color.Green("INFO	Config file loaded at %s", *configPtr)
	}

	color.Yellow("WARN	Standalone mode is %t", *isStandalone)
	// Pre execution checks
	// Check serial scripter connection
	// SSH Server Configureation successful setup
	// Powershell Check

	// Full execution
	logger.InitLogger()

	var wg sync.WaitGroup
	wg.Add(3)

	go heartbeatLoop(*isStandalone)
	go inventoryLoop(*isStandalone)
	go objectMonitoring(*isStandalone)

	wg.Wait()

	select {}

}

func inventoryLoop(standalone bool) {

	if !standalone {
		ticker := time.NewTicker(20 * time.Second)

		for _ = range ticker.C {
			common.PostInventory()
		}
	} else {
		outputName := strconv.FormatInt(time.Now().Unix(), 10) + "_inventory.json"

		inventoryItems := common.GetInventory()

		jsonStr, err := json.Marshal(inventoryItems)
		if err != nil {
			zap.S().Error(err)
			color.Red("JSON marshall error: %v", err)
			os.Exit(0)
		}

		file, err := os.Create(outputName)
		if err != nil {
			zap.S().Error(err)
			color.Red("File creation error: %v", err)
			os.Exit(0)
		}

		_, err = file.WriteString(string(jsonStr))
		if err != nil {
			zap.S().Error(err)
			color.Red("File write Error: %v", err)
			os.Exit(0)
		}

		color.Green("INFO Inventory executed successfully! Output file: %s", outputName)
	}

}

func heartbeatLoop(standalone bool) {

	if !standalone {
		ticker := time.NewTicker(1 * time.Minute)

		for _ = range ticker.C {
			common.HeartBeat()
		}
	} else {
		color.Yellow("INFO	Object Monitoring is not supported in stanalone mode ")
	}

}

func objectMonitoring(standalone bool) {

	if !standalone {
		ticker := time.NewTicker(10 * time.Second)

		for _ = range ticker.C {
			honeymonitor.CreateDirMonitor(common.GetHoneyPaths())
		}
	} else {
		color.Yellow("INFO	Object Monitoring is not supported in stanalone mode ")
	}

}
