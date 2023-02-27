package config

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

type Configuration struct {
	Apis struct {
		Kaspersky struct {
			APIKey string `json:"api_key"`
		} `json:"kaspersky"`
		SerialScripter struct {
			APIKey    string `json:"api_key"`
			UserAgent string `json:"user_agent"`
			URL       string `json:"url"`
		} `json:"serial_scripter"`
	} `json:"apis"`
	ExePaths struct {
		Yara string `json:"yara"`
		Pd   string `json:"pd"`
	} `json:"exe_paths"`
	Honeypaths struct {
		Paths []string `json:"paths"`
	} `json:"honeypaths"`
	Sessions struct {
		whitelist []string `json:"whitelist"`
	} `json:"sessions"`
}

const CONFIG_LOC string = "C:\\Users\\hunte\\Documents\\repos\\wingoEDR\\config.json"

//const CONFIG_LOC string = "C:\\Users\\Administrator\\Downloads\\config.json"

func GetKaperskyKey() string {
	file, _ := os.Open(CONFIG_LOC)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		zap.S().Error("error:", err)
	}
	return configuration.Apis.Kaspersky.APIKey
}

func GetSerialScripterUserAgent() string {
	file, _ := os.Open(CONFIG_LOC)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		zap.S().Error("error:", err)
	}
	return configuration.Apis.SerialScripter.UserAgent
}

func GetYaraExePath() string {
	file, _ := os.Open(CONFIG_LOC)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		zap.S().Error("error:", err)
	}
	return configuration.ExePaths.Yara
}

func GetHoneyPaths() []string {
	file, _ := os.Open(CONFIG_LOC)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		zap.S().Error("error:", err)
	}
	return configuration.Honeypaths.Paths
}

func GetSerialScripterURL() string {
	file, _ := os.Open(CONFIG_LOC)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		zap.S().Error("error:", err)
	}
	return configuration.Apis.SerialScripter.URL
}

func GetWhitelistedUsers() []string {
	file, _ := os.Open(CONFIG_LOC)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		zap.S().Error("error:", err)
	}
	return configuration.Sessions.whitelist
}
