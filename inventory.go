package main

import (
	"runtime"
	"strings"
	"wingoEDR/common"
	"wingoEDR/shares"
)

type InventoryObject struct {
	Name     string        `json:"name"`
	IP       string        `json:"ip"`
	Os       string        `json:"OS"`
	Services []Services    `json:"services"`
	IsOn     bool          `json:"isOn"`
	Docker   []interface{} `json:"docker"`
	Tasks    []interface{} `json:"tasks"`
	Firewall []interface{} `json:"firewall"`
	Shares   []Share       `json:"shares"`
}

type Services struct {
	Port    uint   `json:"port"`
	Service string `json:"service"`
}

type shareAttributes struct {
	Name        string
	Path        string
	Permissions string
}
type SharePermissions struct {
}

func GetInventory() InventoryObject {

	lastOctets := strings.Split(common.GetIP(), ".", 2)
	serialScripterHostName := "host-" + lastOctets[2]

	inv := InventoryObject{
		Name:     serialScripterHostName,
		IP:       common.GetIP(),
		Os:       runtime.GOOS,
		Services: nil,
		IsOn:     true,
		Docker:   nil,
		Tasks:    nil,
		Firewall: nil,
		Shares:   shares.ListSharesWMI(),
	}

	return inv
}
