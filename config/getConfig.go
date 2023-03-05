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
			URL       string `json:"url"`
			UserAgent string `json:"user_agent"`
		} `json:"serial_scripter"`
	} `json:"apis"`
	Blacklist struct {
		Ips []any `json:"ips"`
	} `json:"blacklist"`
	Chainsaw struct {
		Mapping struct {
			Path string `json:"path"`
		} `json:"mapping"`
		Rules struct {
			Path struct {
				Bad    string `json:"bad"`
				Events string `json:"events"`
			} `json:"path"`
		} `json:"rules"`
	} `json:"chainsaw"`
	ExePaths struct {
		Chainsaw string `json:"chainsaw"`
		Yara     string `json:"yara"`
	} `json:"exe_paths"`
	Honeypaths struct {
		Paths []string `json:"paths"`
	} `json:"honeypaths"`
	Whitelist struct {
		Ips      []string `json:"ips"`
		Sessions []string `json:"sessions"`
		Users    []string `json:"users"`
	} `json:"whitelist"`
}

//const CONFIG_LOC string = "C:\\Users\\hunte\\Documents\\repos\\wingoEDR\\config.json"

//const CONFIG_LOC string = "C:\\Users\\Administrator\\Downloads\\config.json"

var (
	CONFIG_LOC string
)

func InitializeConfigLoc(configPath string) {
	CONFIG_LOC = configPath
}

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
	return configuration.Whitelist.Users
}

func GetChainsawPath() string {
	file, _ := os.Open(CONFIG_LOC)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		zap.S().Error("error:", err)
	}
	return configuration.ExePaths.Chainsaw
}

func GetChainsawMapping() string {
	file, _ := os.Open(CONFIG_LOC)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		zap.S().Error("error:", err)
	}
	return configuration.Chainsaw.Mapping.Path
}

func GetChainSawRulesBad() string {
	file, _ := os.Open(CONFIG_LOC)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		zap.S().Error("error:", err)
	}
	return configuration.Chainsaw.Rules.Path.Bad
}
