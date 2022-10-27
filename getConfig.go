package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Apis struct {
		Kaspersky struct {
			APIKey string `json:"api_key"`
		} `json:"kaspersky"`
		SerialScripter struct {
			APIKey    string `json:"api_key"`
			UserAgent string `json:"user_agent"`
		} `json:"serial_scripter"`
	} `json:"apis"`
}

func getKaperskyKey() string {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration.Apis.Kaspersky.APIKey
}

func getSerialScripterUserAgent() string {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration.Apis.SerialScripter.UserAgent
}
