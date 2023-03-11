package main

import (
	"flag"
	"time"
	"wingoEDR/common"
	"wingoEDR/config"
	"wingoEDR/honeymonitor"
	"wingoEDR/logger"
	"wingoEDR/modes"
	"wingoEDR/monitors"
	"wingoEDR/serialscripter"

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
																				 
									Version 0.1.1
									By: Hunter Pittman and the Keyboard Cowboys						 
																				 
	
	`

	println(headline)

	// processIsAdmin, err := common.ProcessIsAdmin()
	// if err != nil {
	// 	zap.S().Fatal("Failed to determine if process is running as admin! Err: %v", err)
	// }

	// if !processIsAdmin {
	// 	zap.S().Fatal("This program must be run as administrator!")
	// }

	defaultConfigPath := config.GenerateConfig()

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

	// var wg sync.WaitGroup
	// wg.Add(1)
	// // chainsawMonitor()

	// // Serial Scripter routines
	// go heartbeatLoop()
	// go inventoryLoop()

	// //Internal routines

	// //go objectMonitoring()

	// wg.Wait()

	// select {}

	monitors.MonitorUsers()
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

func objectMonitoring(standalone bool) {
	ticker := time.NewTicker(10 * time.Second)

	for _ = range ticker.C {
		honeymonitor.CreateDirMonitor(config.GetHoneyPaths())
	}

}

func chainsawMonitor() {
	// ticker := time.NewTicker(10 * time.Second)

	// for _ = range ticker.C {
	// 	chainsaw.FullEventCheck()
	// }

	//chainsaw.RangedEventCheck("2023-03-09T00:00:00", "2023-03-09T023:59:59")
	monitors.FullEventCheck()
}
