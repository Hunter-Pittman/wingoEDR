package common

import (
	"fmt"
	"os"
	"wingoEDR/firewall"
	"wingoEDR/processes"
	"wingoEDR/servicemanager"
	"wingoEDR/shares"
	"wingoEDR/usermanagement"
)

type InventoryObject struct {
	SerialScripterName string                          `json:"name"`
	HostName           string                          `json:"hostname"`
	IP                 string                          `json:"ip"`
	Os                 string                          `json:"OS"`
	Services           []servicemanager.WindowsService `json:"services"`
	Tasks              []interface{}                   `json:"tasks"`
	Firewall           []firewall.FirewallList         `json:"firewall"`
	Shares             []shares.ShareAttributes        `json:"shares"`
	Users              []usermanagement.LocalUser      `json:"users"`
	Processes          []processes.ProcessInfo         `json:"processes"`
}

type InventorySummary struct {
	Name string
	IP   string
	Os   string
}

func GetInventory() InventoryObject {

	processes, _ := processes.GetAllProcesses()

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}

	inv := InventoryObject{
		SerialScripterName: GetSerialScripterHostName(),
		HostName:           hostname,
		IP:                 GetIP(),
		Os:                 OSversion(),
		Services:           servicemanager.Servicelister(),
		Tasks:              nil,
		Firewall:           firewall.FirewallLister(),
		Shares:             shares.ListSharesWMI(),
		Users:              usermanagement.ReturnUsers(),
		Processes:          processes,
	}

	return inv
}

func GetInventorySummary() {

}
