package registrycapture

import (
	"time"

	"go.uber.org/zap"
	"golang.org/x/sys/windows/registry"
)

type RegistryValues struct {
	registryFullPath string
	registryValue    string
}


// supply registry path & the keys you want from those paths
// slices must be of equal length
// EX: GetRegistryValues([]string{"SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion"}, []string{"SystemRoot"})
func GetRegistryValue(registryPathSlice []string, registryKeySlice []string) []RegistryValues {
	var regValues []RegistryValues
	if len(registryPathSlice) != len(registryKeySlice) {
		zap.S().Error("Registry path slice isn't equal to registry key slice!")
	}

	for index, registryPath := range registryPathSlice {
		keyHandle, err := registry.OpenKey(registry.LOCAL_MACHINE, registryPath, registry.QUERY_VALUE)
		if err != nil {
			//logger.Warn("Cannot open registry key:" + registryPath + err.Error())
			zap.S().Warn("Cannot open registry key: " + registryPath + err.Error())
		}
		keyValue, _, err := keyHandle.GetStringValue(registryKeySlice[index])
		if err != nil {
			//logger.Warn("Cannot get value for key: " + registryKeySlice[index])
			zap.S().Warn("Cannot get value for key: " + registryKeySlice[index])
		} else {

			keyValuePair := RegistryValues{registryPath, keyValue}
			regValues = append(regValues, keyValuePair)
		}
		keyHandle.Close()
	}
	return regValues
}

type InstalledSoftware struct {
	Name    string
	Version string
	InstallPath	string
	Publisher	string
	UninstallString	string
}

// Return subkey values for installed software
func GetSoftwareSubkeys(registryPath string) []InstalledSoftware {

	keyHandle, err := registry.OpenKey(registry.LOCAL_MACHINE, registryPath, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		zap.S().Error(err)
	}
	subkeys, err := keyHandle.ReadSubKeyNames(-1)
	if err != nil {
		zap.S().Error(err)
	}
	
	softwareList := []InstalledSoftware{}
	for _, subkey := range subkeys {
		subKeyHandle, err := registry.OpenKey(registry.LOCAL_MACHINE, registryPath +`\`+subkey, registry.QUERY_VALUE)
		if err != nil {
			zap.S().Error(err)
		}
		publisher, _, err := subKeyHandle.GetStringValue("Publisher")
		if err != nil {
			zap.S().Error(err)
		}
		installLocation, _, err := subKeyHandle.GetStringValue("InstallLocation")
		if err != nil {
			zap.S().Error(err)
		}
		uninstallString, _, err := subKeyHandle.GetStringValue("UninstallString")
		if err != nil {
			zap.S().Error(err)
		}
		displayName, _, err := subKeyHandle.GetStringValue("DisplayName")
		if err != nil {
			zap.S().Error(err)
		}
		version, _, err := subKeyHandle.GetStringValue("DisplayVersion")
		if err != nil {
			zap.S().Error(err)
		}
		softwareInfo := InstalledSoftware{displayName, version, installLocation, publisher, uninstallString}
		softwareList = append(softwareList, softwareInfo)
	}
	keyHandle.Close()
	return softwareList
}

func GetLastWriteTime(keyHandle registry.Key) time.Time {
	keyInfo, _ := keyHandle.Stat()
	return keyInfo.ModTime()
}
