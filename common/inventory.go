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
	Services  []servicemanager.WindowsService `json:"services"`
	IsOn      bool                            `json:"isOn"`
	Docker    []interface{}                   `json:"docker"`
	Tasks     []interface{}                   `json:"tasks"`
	Firewall  []firewall.FirewallList         `json:"firewall"`
	Shares    []shares.ShareAttributes        `json:"shares"`
	Processes []processes.ProcessInfo         `json:"processes"`
}

func GetInventory() InventoryObject {

	processes, _ := processes.GetAllProcesses()

	inv := InventoryObject{
		Name:      GetSerialScripterHostName(),
		IP:        GetIP(),
		Os:        runtime.GOOS,
		Services:  servicemanager.Servicelister(), // Maybe move this
		Tasks:     nil,
		Firewall:  firewall.FirewallLister(), // Maybe seperate this
		Shares:    shares.ListSharesWMI(),
		Processes: processes,
	}

	return inv
}
