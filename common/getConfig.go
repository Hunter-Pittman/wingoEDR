package common

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
	ExePaths struct {
		Yara string `json:"yara"`
		Pd   string `json:"pd"`
	} `json:"exe_paths"`
}

func GetKaperskyKey() string {
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

func GetSerialScripterUserAgent() string {
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

func GetYaraExePath() string {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration.ExePaths.Yara
}
