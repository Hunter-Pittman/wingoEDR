package config

import (
	"os"
	"path/filepath"

	"github.com/Jeffail/gabs"
	"github.com/fatih/color"
	"go.uber.org/zap"
)

var (
	externalresourcesPath = "\\externalresources\\"
	configPath            = "\\config.json"
)

func GenerateConfig() string {
	wingoPath, err := os.Executable()
	if err != nil {
		zap.S().Error("Error: ", err.Error())
		color.Red("[ERROR]	An error has been encounterd: ", err.Error())
	}

	wingoFolder := filepath.Dir(wingoPath)
	if err != nil {
		zap.S().Error("Error: ", err.Error())
		color.Red("[ERROR]	An error has been encounterd: ", err.Error())
	}

	fullExternalResourcesPath := wingoFolder + externalresourcesPath
	fullConfigPath := wingoFolder + configPath

	_, err1 := os.Stat(fullExternalResourcesPath)
	if err1 == nil {
		color.Green("[INFO]	External resources folder already exists")
	}
	if os.IsNotExist(err1) {
		color.Red("[ERROR]	External resources folder does not exist, download the external resources and rexecute the program")
		os.Exit(1)
	}

	externalresourcesPath = fullExternalResourcesPath
	configPath = fullConfigPath

	_, err2 := os.Stat(fullConfigPath)
	if err2 == nil {
		color.Green("[INFO]	Config.json already exists, continuing execution...")
		return configPath
	}
	if os.IsNotExist(err2) {
		color.Red("[INFO]	config.json does not exist, generating new one...")
		generateJSON()
	}

	return configPath

}

func generateJSON() {
	jsonOBJ := gabs.New()

	// API Section
	jsonOBJ.Set("no_key", "apis", "kaspersky", "api_key")
	jsonOBJ.Set("random_placeholder", "apis", "serial_scripter", "api_key")
	jsonOBJ.Set("secret", "apis", "serial_scripter", "user_agent")
	jsonOBJ.Set("no_url", "apis", "serial_scripter", "url")

	// EXE Paths
	jsonOBJ.Set(externalresourcesPath+"yara\\yara.exe", "exe_paths", "yara")
	jsonOBJ.Set(externalresourcesPath+"chainsaw\\chainsaw_x86_64-pc-windows-msvc.exe", "exe_paths", "chainsaw")

	// Honeypaths
	jsonOBJ.Array("honeypaths", "paths")

	// Whitelist
	jsonOBJ.Array("whitelist", "ips")
	jsonOBJ.Array("whitelist", "sessions")
	jsonOBJ.Array("whitelist", "users")

	// Blacklist
	jsonOBJ.Array("blacklist", "ips")

	finalJSONObj := jsonOBJ.String()

	file, err := os.Create(configPath)
	if err != nil {
		zap.S().Error("Error: ", err.Error())
		color.Red("[ERROR]	An error has been encounterd: ", err.Error())
	}

	_, err = file.WriteString(finalJSONObj)
	if err != nil {
		zap.S().Error("Error: ", err.Error())
		color.Red("[ERROR]	An error has been encounterd: ", err.Error())
	}

}
