package servicemanager

import (
	"fmt"
	//windows api from github
	wapi "github.com/iamacarpet/go-win64api"
)

type WindowService struct {
	scn string
	dn  string
	st  string
	as  bool
	rp  uint32
}

func servicelister() {
	servslice := make([]WindowService, 0)
	svc, err := wapi.GetServices()
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
	}

	for _, v := range svc {

		helium := WindowService{
			scn: v.SCName,
			dn:  v.DisplayName,
			st:  v.StatusText,
			as:  v.AcceptStop,
			rp:  v.RunningPid}

		servslice = append(servslice, helium)

		fmt.Printf("%-50s - %-75s - Status: %-20s - Accept Stop: %-5t, Running Pid: %d\r\n", v.SCName, v.DisplayName, v.StatusText, v.AcceptStop, v.RunningPid)
	}
	return servslice

}
