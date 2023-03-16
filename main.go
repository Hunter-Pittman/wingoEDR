package main

import (
	"flag"
	"sync"
	"time"
	"wingoEDR/apis/serialscripter"
	"wingoEDR/common"
	"wingoEDR/config"
	"wingoEDR/db"
	"wingoEDR/logger"
	"wingoEDR/modes"
	"wingoEDR/monitors"

	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	// Command line args

	headline := `

	$$\      $$\ $$\            $$$$$$\            $$$$$$$$\ $$$$$$$\  $$$$$$$\  
	$$ | $\  $$ |\__|          $$  __$$\           $$  _____|$$  __$$\ $$  __$$\ 
	$$ |$$$\ $$ |$$\ $$$$$$$\  $$ /  \__| $$$$$$\  $$ |      $$ |  $$ |$$ |  $$ |
	$$ $$ $$\$$ |$$ |$$  __$$\ $$ |$$$$\ $$  __$$\ $$$$$\    $$ |  $$ |$$$$$$$  |
	$$$$  _$$$$ |$$ |$$ |  $$ |$$ |\_$$ |$$ /  $$ |$$  __|   $$ |  $$ |$$  __$$< 
	$$$  / \$$$ |$$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |      $$ |  $$ |$$ |  $$ |
	$$  /   \$$ |$$ |$$ |  $$ |\$$$$$$  |\$$$$$$  |$$$$$$$$\ $$$$$$$  |$$ |  $$ |
	\__/     \__|\__|\__|  \__| \______/  \______/ \________|\_______/ \__|  \__|
																				 
									Version 0.1.2
									By: Hunter Pittman and the Keyboard Cowboys						 
																				 
	
	`

	println(headline)

	// Initializations
	db.DbInit()
	defaultConfigPath := config.GenerateConfig()
	runRequest := serialscripter.CheckEndpoint()
	if !runRequest {
		zap.S().Warn("Serial Scripter is not running. WingoEDR will continue to run in offline mode.")
	}

	configPtr := flag.String("config", defaultConfigPath, "Provide path to the config file")
	mode := flag.String("mode", "default", "List what mode you would like wingoEDR to execute in. The default is to enable continous monitoring.")

	// Backup flags
	backupDir := flag.String("backupdir", "C:\\backups", "Enter the path where your backups are going to be stored.")
	backupItem := flag.String("backupitem", "", "Enter the path to the file or directory you wish to backup.")

	// Decompress flags
	decompressItem := flag.String("decompressitem", "", "Enter the path to the file or directory you wish to decompress")

	// Chainsaw flags
	from := flag.String("from", "", "Enter the start timestamp in the format of YYYY-MM-DDTHH:MM:SS")
	to := flag.String("to", "", "Enter the end timestamp in the format of YYYY-MM-DDTHH:MM:SS")
	json := flag.Bool("json", false, "Enter true to output in json format")

	flag.Parse()

	common.VerifyWindowsPathFatal(*configPtr)
	zap.S().Infof("Config file loaded %s", *configPtr)

	config.InitializeConfigLoc(*configPtr)

	// thing := "chainsaw" // TEST VALUE
	// mode = &thing       // TEST VALUE
	// thing2 := true      // TEST VALUE
	// json = &thing2      // TEST VALUE

	modes.ModeHandler(*mode, map[string]modes.Params{"backupDir": *backupDir, "backupItem": *backupItem, "decompressItem": *decompressItem, "from": *from, "to": *to, "json": *json})

	// Pre execution checks
	// Check serial scripter connection
	// SSH Server configuration successful setup
	// Powershell Check

	// continousMonitoring()

	continousMonitoring()

}

func continousMonitoring() {
	var wg sync.WaitGroup
	wg.Add(3)

	// Serial Scripter routines
	// go heartbeatLoop()
	// go inventoryLoop()

	// Internal routines
	//go userLoop()
	//go smbShareLoop()
	// go serviceLoop()
	// go chainsawLoop()
	processLoop()

	wg.Wait()

	select {}

}

func processLoop() {
	// monitors.InitProcesses()
	// ticker := time.NewTicker(1 * time.Minute)

	// for _ = range ticker.C {
	// 	monitors.ProcessMonitor()
	// }

	monitors.ProcessMonitor()
}

func smbShareLoop() {
	monitors.InitShares()
	ticker := time.NewTicker(10 * time.Second)

	for _ = range ticker.C {
		monitors.SharesMonitor()
	}
}

func inventoryLoop() {
	ticker := time.NewTicker(20 * time.Second)

	for _ = range ticker.C {
		serialscripter.PostInventory()
	}

}

func heartbeatLoop() {
	ticker := time.NewTicker(1 * time.Minute)

	for _ = range ticker.C {
		serialscripter.HeartBeat()
	}
}

func userLoop() {
	monitors.InitUsers()
	ticker := time.NewTicker(1 * time.Minute)

	for _ = range ticker.C {
		monitors.MonitorUsers()
	}
}

func chainsawLoop() {
	monitors.FullEventCheck()
	// ticker := time.NewTicker(1 * time.Minute)

	// for _ = range ticker.C {
	// 	chainsaw.FullEventCheck()
	// }

	currentTime := time.Now()
	oneMinuteAgo := currentTime.Add(-1 * time.Minute)
	newTimestamp := oneMinuteAgo.Format("2006-01-02T15:04:05")

	monitors.RangedEventCheck(newTimestamp, currentTime.Format("2006-01-02T15:04:05"))
}
