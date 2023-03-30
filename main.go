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

	"github.com/kardianos/service"
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
																				 
									Version v0.1.3-alpha
									By: Hunter Pittman and the Keyboard Cowboys						 
																				 
	
	`

	println(headline)

	// Initializations
	db.DbInit()
	defaultConfigPath := config.GenerateConfig()

	configPtr := flag.String("config", defaultConfigPath, "Provide path to the config file")
	mode := flag.String("mode", "default", "List what mode you would like wingoEDR to execute in. The default is to enable continous monitoring.")
	offline := flag.Bool("offline", false, "Use this flag to diasble posting to SerialScripter.")

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
	runRequest := serialscripter.CheckEndpoint(*offline)
	if !runRequest {
		zap.S().Warn("Serial Scripter is not running. WingoEDR will continue to run in offline mode.")
	}

	// thing := "decompress"                         // TEST VALUE
	// mode = &thing                                 // TEST VALUE
	// thing2 := `C:\backups\compressed_procexp.exe` // TEST VALUE
	// decompressItem = &thing2                      // TEST VALUE

	paramItems := map[string]modes.Params{"backupDir": *backupDir, "backupItem": *backupItem, "decompressItem": *decompressItem, "from": *from, "to": *to, "json": *json}
	modes.ModeHandler(*mode, paramItems)

	// Pre execution checks
	// Check serial scripter connection
	// SSH Server configuration successful setup
	// Powershell Check

	// continousMonitoring()

	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}

	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		zap.S().Error("Cannot create the service: " + err.Error())
	}
	err = s.Run()
	if err != nil {
		zap.S().Error("Cannot start the service: " + err.Error())
	}

}

func continousMonitoring() {
	var wg sync.WaitGroup
	wg.Add(5)

	// Serial Scripter routines
	//go heartbeatLoop()
	//go inventoryLoop()

	// Internal routines
	//go userLoop()
	go smbShareLoop()
	//go serviceLoop()
	//go chainsawLoop()
	//go processLoop()
	//go autorunsLoop()
	//go softwareLoop()

	wg.Wait()

	select {}

}

func autorunsLoop() {
	monitors.InitAutoruns()
	ticker := time.NewTicker(10 * time.Second)

	for _ = range ticker.C {
		monitors.MonitorAutoruns()
	}
}

func processLoop() {
	// monitors.InitProcesses()
	// ticker := time.NewTicker(1 * time.Minute)

	// for _ = range ticker.C {
	// 	monitors.ProcessMonitor()
	// }

	//monitors.ProcessMonitor()
}

func smbShareLoop() {
	monitors.InitShares()
	ticker := time.NewTicker(10 * time.Second)

	for _ = range ticker.C {
		monitors.SharesMonitor()
	}
}

func softwareLoop() {
	monitors.InitSoftware()
	ticker := time.NewTicker(30 * time.Second)

	for _ = range ticker.C {
		monitors.SoftwareMonitor()
	}
}

func serviceLoop() {
	monitors.InitServices()
	ticker := time.NewTicker(30 * time.Second)

	for _ = range ticker.C {
		monitors.MonitorServices()
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
	monitors.FullEventReportCheck()

	ticker := time.NewTicker(60 * time.Second)
	processendtime := time.Now().Format("2006-01-02T15:04:05")
	for _ = range ticker.C {
		currentTime := time.Now()
		// oneMinuteAgo := currentTime.Add(-72 * time.Second)
		// newTimestamp := oneMinuteAgo.Format("2006-01-02T15:04:05")
		monitors.RangedEventReportCheck(processendtime, currentTime.Format("2006-01-02T15:04:05"))
		processendtime = time.Now().Format("2006-01-02T15:04:05")
	}

}

// SERVICE JUNK
const serviceName = "wingoEDR"
const serviceDescription = "Bleh"

type program struct{}

func (p program) Start(s service.Service) error {
	zap.S().Info(s.String() + " started")
	go p.run()
	return nil
}

func (p program) Stop(s service.Service) error {
	zap.S().Info(s.String() + " stopped")
	return nil
}

func (p program) run() {
	continousMonitoring()
}
