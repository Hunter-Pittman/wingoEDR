package inventory

import (
	"os"
	"time"
	"wingoEDR/common"
	"wingoEDR/firewall"
	"wingoEDR/processes"
	"wingoEDR/registrycapture"
	"wingoEDR/servicemanager"
	"wingoEDR/shares"
	"wingoEDR/usermanagement"

	"go.uber.org/zap"
)

type InventoryObject struct {
	SerialScripterName string                              `json:"name"`
	HostName           string                              `json:"hostname"`
	IP                 string                              `json:"ip"`
	Os                 string                              `json:"OS"`
	Services           []servicemanager.WindowsService     `json:"services"`
	Tasks              []interface{}                       `json:"tasks"`
	Firewall           []firewall.FirewallList             `json:"FirewallList"`
	InstalledSoftware  []registrycapture.InstalledSoftware `json:"installedSoftware"`
	Shares             []shares.SMBInfo                    `json:"shares"`
	Users              []usermanagement.User               `json:"users"`
	Processes          []processes.ProcessInfo             `json:"processes"`
	TimeConnected      string                              `json:"timeConnected"`
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
		InstalledSoftware:  registrycapture.GetSoftwareSubkeys(`HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\Current Version\Uninstall`),
		Shares:             shares.GetShares(),
		Users:              usermanagement.ReturnUsers(),
		Processes:          processes,
		TimeConnected:      time.Now().Format("03:04:05 PM"),
	}

	return inv
}

func GetInventorySummary() {

}
