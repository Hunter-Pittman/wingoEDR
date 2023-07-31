package inventory

import (
	"os"
	"time"
	"wingoEDR/common"
	"wingoEDR/firewall"
	"wingoEDR/installedsoftware"
	"wingoEDR/processes"
	"wingoEDR/servicemanager"
	"wingoEDR/shares"
	"wingoEDR/usermanagement"

	"go.uber.org/zap"
)

type InventoryObject struct {
	SerialScripterName string                            `json:"name"`
	HostName           string                            `json:"hostname"`
	IP                 string                            `json:"ip"`
	Os                 string                            `json:"OS"`
	Services           []servicemanager.WindowsService   `json:"services"`
	Tasks              []interface{}                     `json:"tasks"`
	Firewall           []firewall.FirewallList           `json:"FirewallList"`
	InstalledSoftware  []installedsoftware.Win32_Product `json:"installedSoftware"`
	Shares             []shares.SMBInfo                  `json:"shares"`
	Users              []usermanagement.User             `json:"users"`
	Processes          []processes.ProcessInfo           `json:"processes"`
	TimeConnected      string                            `json:"timeConnected"`
}

func GetInventory() InventoryObject {

	processes, _ := processes.GetAllProcesses()

	hostname, err := os.Hostname()
	if err != nil {
		zap.S().Error(err)
	}

	inv := InventoryObject{
		SerialScripterName: common.GetSerialScripterHostName(),
		HostName:           hostname,
		IP:                 common.GetIP(),
		Os:                 common.OSversion(),
		Services:           servicemanager.Servicelister(),
		Tasks:              nil,
		Firewall:           firewall.FirewallLister(),
		InstalledSoftware:  installedsoftware.WmiSoftwareQuery(),
		Shares:             shares.GetShares(),
		Users:              usermanagement.ReturnUsers(),
		Processes:          processes,
		TimeConnected:      time.Now().Format("03:04:05 PM"),
	}

	return inv
}
