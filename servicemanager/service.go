package servicemanager

import (
	"fmt"
	//windows api from github
	wapi "github.com/iamacarpet/go-win64api"
)

func servicelister() {

	// sets the service and error varibles

	svc, err := wapi.GetServices()
	//checks for an error in the code
	if err != nil {
		//if there is one it prints it out
		fmt.Print("%s\r\n", err.Error())
	}

	//enumerates through the list of services
	for _, v := range svc {
		//prints each service with varibles attached
		fmt.Printf("%-50s - %-75s - Status: %-20s - Accept Stop: %-5t, Running Pid: %d\r\n", v.SCName, v.DisplayName, v.StatusText, v.AcceptStop, v.RunningPid)
	}

}
