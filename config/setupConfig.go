package config

import (
	"fmt"
	"os"
	"path/filepath"
	"wingoEDR/unzip"

	"github.com/Jeffail/gabs"
	"github.com/cavaliergopher/grab/v3"
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
	}

	wingoFolder := filepath.Dir(wingoPath)
	if err != nil {
		zap.S().Error("Error: ", err.Error())
	}

	fullExternalResourcesPath := wingoFolder + externalresourcesPath
	fullConfigPath := wingoFolder + configPath

	_, err1 := os.Stat(fullExternalResourcesPath)
	if err1 == nil {
		zap.S().Info("External resources folder already exists")
	}
	if os.IsNotExist(err1) {
		zipFolder := wingoFolder + "\\externalresources.zip"
		_, err := os.Stat(zipFolder)
		if err == nil {
			zap.S().Info("External resources zip file already exists, extracting...")
			unzip.Unzip(zipFolder)
		} else {
			zap.S().Warn("External resources folder does not exist")
			releaseVersion := "v0.1.3-alpha"
			zap.S().Warnf("Attempting download of external resources %s...", releaseVersion)
			fullUrl := fmt.Sprintf("https://github.com/Hunter-Pittman/wingoEDR/releases/download/%s/externalresources.zip", releaseVersion)
			_, err := grab.Get(wingoFolder, fullUrl)
			if err != nil {
				zap.S().Fatal("Unable to download external resources: ", err)
			}
			unzip.Unzip(zipFolder)
			//os.Exit(1)
		}

	}

	externalresourcesPath = fullExternalResourcesPath
	configPath = fullConfigPath

	_, err2 := os.Stat(fullConfigPath)
	if err2 == nil {
		zap.S().Info("Config.json already exists, continuing execution...")
		return configPath
	}
	if os.IsNotExist(err2) {
		zap.S().Info("config.json does not exist, generating new one...")
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
	jsonOBJ.Set("no_key", "apis", "siem", "api_key")
	jsonOBJ.Set("no_url", "apis", "siem", "url")

	// EXE Paths
	jsonOBJ.Set(externalresourcesPath+"yara\\yara.exe", "exe_paths", "yara")
	jsonOBJ.Set(externalresourcesPath+"chainsaw\\chainsaw_x86_64-pc-windows-msvc.exe", "exe_paths", "chainsaw")

	// Chainsaw
	jsonOBJ.Set(externalresourcesPath+"chainsaw\\rules\\Bad\\", "chainsaw", "rules", "path", "bad")
	jsonOBJ.Set(externalresourcesPath+"chainsaw\\rules\\Events\\", "chainsaw", "rules", "path", "events")
	jsonOBJ.Set(externalresourcesPath+"chainsaw\\mappings\\sigma-event-logs-all.yml", "chainsaw", "mapping", "path")

	// Honeypaths
	jsonOBJ.Array("honeypaths", "paths")

	// Blacklist
	jsonOBJ.Array("lists", "blacklist", "ips")
	jsonOBJ.Array("lists", "blacklist", "sessions")
	jsonOBJ.Array("lists", "blacklist", "users")
	jsonOBJ.Array("lists", "blacklist", "publishers")
	jsonOBJ.Array("lists", "blacklist", "processes")

	// Graylist
	jsonOBJ.Array("lists", "graylist", "ips")
	jsonOBJ.Array("lists", "graylist", "sessions")
	jsonOBJ.Array("lists", "graylist", "users")
	jsonOBJ.Array("lists", "graylist", "publishers")
	jsonOBJ.Array("lists", "graylist", "processes")

	// Whitelist
	jsonOBJ.Array("lists", "whitelist", "ips")
	jsonOBJ.Array("lists", "whitelist", "sessions")
	jsonOBJ.Array("lists", "whitelist", "users")
	jsonOBJ.Array("lists", "whitelist", "publishers")
	jsonOBJ.Array("lists", "whitelist", "processes")

	finalJSONObj := jsonOBJ.String()

	file, err := os.Create(configPath)
	if err != nil {
		zap.S().Error("Error: ", err.Error())
	}

	_, err = file.WriteString(finalJSONObj)
	if err != nil {
		zap.S().Error("Error: ", err.Error())
	}

}
