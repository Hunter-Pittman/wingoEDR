package common

import (
	"runtime"
	"wingoEDR/firewall"
	"wingoEDR/processes"
	"wingoEDR/servicemanager"
	"wingoEDR/shares"
)

type InventoryObject struct {
	Name      string                          `json:"name"`
	IP        string                          `json:"ip"`
	Os        string                          `json:"OS"`
	Services  []servicemanager.WindowsService `json:"windowsservices"`
	Tasks     []interface{}                   `json:"tasks"`
	Firewall  []firewall.FirewallList         `json:"firewall"`
	Shares    []shares.ShareAttributes        `json:"shares"`
	Processes []processes.ProcessInfo         `json:"processes"`
}

type InventorySummary struct {
	Name string
	IP   string
	Os   string
}

func GetInventory() InventoryObject {

	processes, _ := processes.GetAllProcesses()

	inv := InventoryObject{
		Name:      GetSerialScripterHostName(),
		IP:        GetIP(),
		Os:        runtime.GOOS,
		Services:  servicemanager.Servicelister(),
		Tasks:     nil,
		Firewall:  firewall.FirewallLister(),
		Shares:    shares.ListSharesWMI(),
		Processes: processes,
	}

	return inv
}

func GetInventorySummary() {

}
