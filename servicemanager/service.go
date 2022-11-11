package servicemanager

import (
	//windows api from github
	wapi "github.com/iamacarpet/go-win64api"
	"go.uber.org/zap"
)

type WindowsService struct {
	SCname      string
	DisplayName string
	StatusText  string
	AcceptStop  bool
	RunningPID  uint32
}

func Servicelister() []WindowsService {
	servslice := make([]WindowsService, 0)
	svc, err := wapi.GetServices()
	if err != nil {
		zap.S().Error("Getting services failed!")
	}

	for _, v := range svc {

		helium := WindowsService{
			SCname:      v.SCName,
			DisplayName: v.DisplayName,
			StatusText:  v.StatusText,
			AcceptStop:  v.AcceptStop,
			RunningPID:  v.RunningPid}

		servslice = append(servslice, helium)

	}
	return servslice
}
