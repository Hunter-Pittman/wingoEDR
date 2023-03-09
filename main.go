package main

import (
	"flag"
	"sync"
	"time"
	"wingoEDR/chainsaw"
	"wingoEDR/common"
	"wingoEDR/config"
	"wingoEDR/honeymonitor"
	"wingoEDR/logger"
	"wingoEDR/modes"

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
																				 
									Version 0.1.0
									By: Hunter Pittman and the Keyboard Cowboys						 
																				 
	
	`

	println(headline)

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

	flag.Parse()

	common.VerifyWindowsPathFatal(*configPtr)
	//color.Green("[INFO]	Config file loaded %s", *configPtr)
	zap.S().Infof("Config file loaded %s", *configPtr)

	config.InitializeConfigLoc(*configPtr)

	// thing := "chainsaw" // TEST VALUE
	// mode = &thing       // TEST VALUE

	modes.ModeHandler(*mode, map[string]string{"backupDir": *backupDir, "backupItem": *backupItem, "decompressItem": *decompressItem, "from": *from, "to": *to})

	// Pre execution checks
	// Check serial scripter connection
	// SSH Server configuration successful setup
	// Powershell Check

	var wg sync.WaitGroup
	wg.Add(2)
	//chainsawMonitor()

	// // Serial Scripter routines
	go heartbeatLoop()
	go inventoryLoop()

	// //Internal routines

	// go objectMonitoring(*isStandalone)

	wg.Wait()

	select {}

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

func objectMonitoring(standalone bool) {
	ticker := time.NewTicker(10 * time.Second)

	for _ = range ticker.C {
		honeymonitor.CreateDirMonitor(config.GetHoneyPaths())
	}

}

func chainsawMonitor() {
	chainsaw.FullEventCheck()
	// ticker := time.NewTicker(1 * time.Minute)

	// for _ = range ticker.C {

	// }
}
