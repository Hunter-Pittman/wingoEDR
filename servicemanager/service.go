package servicemanager

import (
	//windows api from github
	wapi "github.com/iamacarpet/go-win64api"
	"go.uber.org/zap"
)

type WindowService struct {
	SCname      string
	displayName string
	statusText  string
	acceptStop  bool
	runningPID  uint32
}

func Servicelister() []WindowService {
	servslice := make([]WindowService, 0)
	svc, err := wapi.GetServices()
	if err != nil {
		zap.S().Error("Getting services failed!")
	}

	for _, v := range svc {

		helium := WindowService{
			SCname:      v.SCName,
			displayName: v.DisplayName,
			statusText:  v.StatusText,
			acceptStop:  v.AcceptStop,
			runningPID:  v.RunningPid}

		servslice = append(servslice, helium)

	}
	return servslice

}
