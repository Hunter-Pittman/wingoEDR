package main

import (
	"wingoEDR/common"
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

type Share struct {
	Name        string
	Fullpath    string
	Permissions string
}

type SharePermissions struct {
}

func GetInventory() InventoryObject {
	var services []Services

	services = append(services, Services{
		Port:    22,
		Service: "SSH",
	})

	var shares []Share

	shares = append(shares, Share{
		Name:        "My share",
		Fullpath:    "C:\\Users\\hunte",
		Permissions: "My permissions",
	})

	//newName := strings.SplitAfterN(3)

	temp := InventoryObject{
		Name:     "host-299",
		IP:       common.GetIP(),
		Os:       "Windows 11",
		Services: services,
		IsOn:     true,
		Docker:   nil,
		Tasks:    nil,
		Firewall: nil,
		Shares:   shares,
	}

	return temp
}
