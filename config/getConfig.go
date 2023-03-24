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
			Url       string `json:"url"`
			UserAgent string `json:"user_agent"`
		} `json:"serial_scripter"`
		Siem struct {
			APIKey string `json:"api_key"`
			Url    string `json:"url"`
		} `json:"siem"`
	} `json:"apis"`
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
		Paths []any `json:"paths"`
	} `json:"honeypaths"`
	Lists struct {
		Blacklist struct {
			Ips        []any `json:"ips"`
			Processes  []any `json:"processes"`
			Publishers []any `json:"publishers"`
			Sessions   []any `json:"sessions"`
			Users      []any `json:"users"`
		} `json:"blacklist"`
		Graylist struct {
			Ips        []any `json:"ips"`
			Processes  []any `json:"processes"`
			Publishers []any `json:"publishers"`
			Sessions   []any `json:"sessions"`
			Users      []any `json:"users"`
		} `json:"graylist"`
		Whitelist struct {
			Ips        []any `json:"ips"`
			Processes  []any `json:"processes"`
			Publishers []any `json:"publishers"`
			Sessions   []any `json:"sessions"`
			Users      []any `json:"users"`
		} `json:"whitelist"`
	} `json:"lists"`
}

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

func GetHoneyPaths() []any {
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

// func GetSerialScripterURL() string {

// 	// Read the JSON data from the file
// 	data, err := ioutil.ReadFile(CONFIG_LOC)
// 	if err != nil {
// 		zap.S().Error("error:", CONFIG_LOC)
// 		zap.S().Error("error:", err)
// 	}

// 	jsonParsed, err := gabs.ParseJSON(data)
// 	if err != nil {
// 		zap.S().Error("error:", err)
// 	}

// 	url, ok := jsonParsed.Path("apis.serial_scripter.url").Data().(string)
// 	if ok == false {
// 		zap.S().Error("key not find:", err)
// 	}

// 	return url
// }

func GetSerialScripterURL() string {
	file, _ := os.Open(CONFIG_LOC)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		zap.S().Error("error:", err)
	}
	return configuration.Apis.SerialScripter.Url
}

func GetWhitelistedUsers() []any {
	file, _ := os.Open(CONFIG_LOC)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		zap.S().Error("error:", err)
	}
	return configuration.Lists.Whitelist.Users
}

func GetGraylistedProcesses() []any {
	file, _ := os.Open(CONFIG_LOC)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		zap.S().Error("error:", err)
	}
	return configuration.Lists.Graylist.Processes
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

func GetSiemApiKey() string {
	file, _ := os.Open(CONFIG_LOC)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		zap.S().Error("error:", err)
	}
	return configuration.Apis.Siem.APIKey
}

func GetSiemUrl() string {
	file, _ := os.Open(CONFIG_LOC)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		zap.S().Error("error:", err)
	}
	return configuration.Apis.Siem.Url
}
