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
		//zap.S().Fatal("Registry path slice isn't equal to registry key slice!")
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

func GetLastWriteTime(keyHandle registry.Key) time.Time {
	keyInfo, _ := keyHandle.Stat()
	return keyInfo.ModTime()
}
