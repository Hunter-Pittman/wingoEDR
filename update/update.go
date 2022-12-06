package update

import (
	"os/exec"
	"fmt"
	"encoding/json"
	"go.uber.org/zap"
)

// Can use https://www.partitionwizard.com/partitionmagic/powershell-windows-update.html to install windows updates
// But requires installing additional Powershell module
// Also requires nuget: https://www.tutorialspoint.com/how-to-install-the-nuget-package-using-powershell
// But both can be installed with commands

type UpdateInfo struct {
	Major	string
	Minor	string
	Build	string
	Revision	string
	MajorRevision	string
	MinorRevision	string
	HotFixDesc	string
	HotFixID	string
	HotFixInstalledBy	string
	HotFixInstallDate	string
}


// Will install wuget & PSWindowsUpdate
// Needs to be tested
// Performs update
func PerformUpdate() {
	// Install package manager & module
	_, _ = exec.Command("powershell.exe", "Install-PackageProvider -Name Nuget -Force").Output()
	_, _ = exec.Command("powershell.exe", "Install-Module -Name PSWindowsUpdate -Force").Output()

	// Run update mechanism
	_, err := exec.Command("powershell.exe", "Import-Module PSWindowsUpdate").Output()
	if err != nil {
		zap.S().Error("Error importing PSWindowsUpdate Module: " + err.Error())
	}
	_, err = exec.Command("powershell.exe", "Get-WindowsUpdate -AcceptAll -Install -AutoReboot").Output()
	if err != nil {
		zap.S().Error("Error executing Get-WindowsUpdate: " + err.Error())
	}

}

func ReturnHotFixInfo() map[string]interface{} {
	unformattedOut, err := exec.Command("powershell.exe", "(Get-HotFix | Sort-Object -Property InstalledOn)[-1] | Select InstalledOn, HotFixID, Source, Description, InstalledBy | ConvertTo-JSON").Output()
	if err != nil {
		zap.S().Error("Error getting hot fix update info: " + err.Error())
	}
	var jsonOut map[string]interface{}
	//json.Unmarshal([]byte(string(unformattedOut)), &jsonOut)
	json.Unmarshal(unformattedOut, &jsonOut)
	return jsonOut
}


// Call this function to return UpdateInfo struct with info about hotfix & OS version
func ReturnOSVersion() UpdateInfo {
	unformattedOut, err := exec.Command("powershell.exe", "[System.Environment]::OSVersion.Version | ConvertTo-JSON").Output()
	if err != nil {
		zap.S().Error("Error getting OS version info: " + err.Error())
	}
	// turn output into readable json format
	var jsonOut map[string]interface{}
	json.Unmarshal(unformattedOut, &jsonOut)
	
	hotFixInfo := ReturnHotFixInfo()


	OSVersionInfo := UpdateInfo{InterfaceToString(jsonOut["Major"]), InterfaceToString(jsonOut["Minor"]), InterfaceToString(jsonOut["Build"]), InterfaceToString(jsonOut["Revision"]), InterfaceToString(jsonOut["MajorRevision"]), InterfaceToString(jsonOut["MinorRevision"]), InterfaceToString(hotFixInfo["Description"]), InterfaceToString(hotFixInfo["HotFixID"]), InterfaceToString(hotFixInfo["InstalledBy"]), InterfaceToString(hotFixInfo["InstalledOn"])}
	return OSVersionInfo	
	
}

func InterfaceToString(data interface{}) string {
	dataStr := fmt.Sprintf("%v", data)
	return dataStr
}