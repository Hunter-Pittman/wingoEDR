package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
	"wingoEDR/autoruns"
	"wingoEDR/backup"
	"wingoEDR/common"
	"wingoEDR/honeymonitor"
	"wingoEDR/logger"

	"github.com/fatih/color"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	// Command line args
	defaultConfigPath := "C:\\Users\\FORENSICS\\AppData\\Roaming\\wingoEDR\\config.json"
	//defaultConfigPath := "C:\\Users\\hunte\\Documents\\repos\\wingoEDR\\config.json"
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

	color.Yellow("[WARN]	Standalone mode is %t", *isStandalone)

	switch *mode {
	case "backup":
		color.Green("[INFO]	Mode is %s", *mode)

		// Required params check
		if *backupItem == "" {
			color.Red("[ERROR]	--backupitem is an essential flag for this mode!", *isStandalone)
			return
		}

		common.VerifyWindowsPathFatal(*backupDir)
		common.VerifyWindowsPathFatal(*backupItem)
		file, err := os.Open(*backupItem)
		if err != nil {
			log.Fatal("Backup item file access failure! Err: %v", err)
		}

		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatal("Backup item file access failure! Err: %v", err)
		}

		if fileInfo.IsDir() {
			backup.BackDir(*backupItem, false)
			color.Green("[INFO]	Backup of %s is complete!", *backupItem)
		} else {
			newFileName := "\\compressed_" + fileInfo.Name()
			backup.BackFile(newFileName, *backupDir)
			color.Green("[INFO]	Backup of %s is complete!", *backupItem)
		}

		return
	case "chainsaw":
		color.Green("[INFO]	Mode is %s", *mode)
	case "sessions":
		color.Green("[INFO]	Mode is %s", *mode)
	case "decompress":
		color.Green("[INFO]	Mode is %s", *mode)
		reader, err := os.Open(*decompressItem)
		if err != nil {
			log.Fatal("Backup item file access failure! Err: %v", err)
		}

		writer, err := os.Create(".\\decompressed_lol.txt")
		if err != nil {
			log.Fatal("Backup item file access failure! Err: %v", err)
		}
		common.Decompress(reader, writer)
		return

	default:
		color.Green("[WARN]	Mode is %s", *mode)

	}

	// Pre execution checks
	// Check serial scripter connection
	// SSH Server Configureation successful setup
	// Powershell Check

	// Full execution

	autoruns.FullAutorunsDump()

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
