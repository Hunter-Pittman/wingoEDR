package main

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"sync"
	"time"
	"wingoEDR/common"
	"wingoEDR/config"
	"wingoEDR/honeymonitor"
	"wingoEDR/logger"
	"wingoEDR/modes"

	"github.com/fatih/color"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	// Command line args

	defaultConfigPath := config.GenerateConfig()

	configPtr := flag.String("config", defaultConfigPath, "Provide path to the config file")
	isStandalone := flag.Bool("standalone", false, "If serial scripter is not available then it outputs datga in local csv")
	mode := flag.String("mode", "default", "List what mode you would like wingoEDR to execute in. The default is to enable continous monitoring.")

	// Backup flags
	backupDir := flag.String("backupdir", "C:\\backups", "Enter the path where your backups are going to be stored.")
	backupItem := flag.String("backupitem", "", "Enter the path to the file or directory you wish to backup.")

	// Decompress flags
	decompressItem := flag.String("decompressitem", "", "Enter the path to the file or directory you wish to decompress")

	flag.Parse()

	common.VerifyWindowsPathFatal(*configPtr)
	color.Green("[INFO]	Config file loaded %s", *configPtr)

	config.InitializeConfigLoc(*configPtr)

	color.Yellow("[WARN]	Standalone mode is %t", *isStandalone)

	thing := "chainsaw"
	mode = &thing

	modes.ModeHandler(*mode, map[string]string{"backupDir": *backupDir, "backupItem": *backupItem, "decompressItem": *decompressItem})

	// Pre execution checks
	// Check serial scripter connection
	// SSH Server Configureation successful setup
	// Powershell Check

	// Full execution

	var wg sync.WaitGroup
	wg.Add(3)

	go heartbeatLoop(*isStandalone)
	go inventoryLoop(*isStandalone)
	//go objectMonitoring(*isStandalone)

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
			honeymonitor.CreateDirMonitor(config.GetHoneyPaths())
		}
	} else {
		color.Yellow("INFO	Object Monitoring is not supported in stanalone mode ")
	}

}
