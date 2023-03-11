package servicemanager

import (
	//windows api from github
	"math"
	"wingoEDR/processes"

	wapi "github.com/iamacarpet/go-win64api"
	"go.uber.org/zap"
)

var procs, _ = processes.GetAllProcesses()

type WindowsService struct {
	SCname      string
	DisplayName string
	StatusText  string
	AcceptStop  bool
	RunningPID  uint32
	Port        uint32 `json:"port"`
}

func Servicelister() []WindowsService {
	servslice := make([]WindowsService, 0)
	svc, err := wapi.GetServices()
	if err != nil {
		zap.S().Error("Getting services failed!")
	}

	for _, v := range svc {

		var servicePort uint32 = math.MaxUint32
		for _, process := range procs {
			if uint32(process.Pid) == v.RunningPid {
				if len(process.NetworkConnections) > 0 {
					servicePort = process.NetworkConnections[0].LocalPort
				}
			}
		}

		if servicePort == math.MaxUint32 {
			servicePort = 0
		}

		helium := WindowsService{
			SCname:      v.SCName,
			DisplayName: v.DisplayName,
			StatusText:  v.StatusText,
			AcceptStop:  v.AcceptStop,
			RunningPID:  v.RunningPid,
			Port:        servicePort,
		}

		servslice = append(servslice, helium)

	}
	return servslice
}
