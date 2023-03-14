package inventory

import (
	"fmt"
	"os"
	"wingoEDR/common"
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
	Shares             []shares.SMBInfo                `json:"shares"`
	Users              []usermanagement.User           `json:"users"`
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
		SerialScripterName: common.GetSerialScripterHostName(),
		HostName:           hostname,
		IP:                 common.GetIP(),
		Os:                 common.OSversion(),
		Services:           servicemanager.Servicelister(),
		Tasks:              nil,
		Firewall:           firewall.FirewallLister(),
		Shares:             shares.GetShares(),
		Users:              usermanagement.ReturnUsers(),
		Processes:          processes,
	}

	return inv
}

func GetInventorySummary() {

}
