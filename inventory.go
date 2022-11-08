package main

import (
	"runtime"
	"strings"
	"wingoEDR/common"
	"wingoEDR/servicemanager"
	"wingoEDR/shares"
)

type InventoryObject struct {
	Name     string                          `json:"name"`
	IP       string                          `json:"ip"`
	Os       string                          `json:"OS"`
	Services []servicemanager.WindowsService `json:"services"`
	IsOn     bool                            `json:"isOn"`
	Docker   []interface{}                   `json:"docker"`
	Tasks    []interface{}                   `json:"tasks"`
	Firewall []interface{}                   `json:"firewall"`
	Shares   []shares.ShareAttributes        `json:"shares"`
}

func GetInventory() InventoryObject {
	lastOctets := strings.Split(common.GetIP(), ".")
	serialScripterHostName := "host-" + lastOctets[3]

	inv := InventoryObject{
		Name:     serialScripterHostName,
		IP:       common.GetIP(),
		Os:       runtime.GOOS,
		Services: servicemanager.Servicelister(),
		Docker:   nil,
		Tasks:    nil,
		Firewall: nil,
		Shares:   shares.ListSharesWMI(),
	}

	return inv
}
