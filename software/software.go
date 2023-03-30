package software

import (
	"encoding/json"
	"os/exec"
)

type Software struct {
	Name    string `json:"Name"`
	Version string `json:"Version"`
	Vendor  string `json:"Vendor"`
}

func GetInstalledSofware() []Software {
	out, _ := exec.Command("powershell.exe", "get-wmiobject -Class Win32_Product |Select-Object Name, Version, Vendor | ConvertTo-JSON").Output()
	softwareList := []Software{}
	//var softwareJson map[string]interface{}
	_ = json.Unmarshal(out, &softwareList)
	// parsing out
	/*(for index, line := range strings.Split(string(out), "\n") {
		// ignore the first & second lines, they're just headers
		if index == 0 || index == 1 {
			continue
		} else {
			tabSplit := strings.Split(line, ",")
			// ensure we're parsing a valid line
			if len(tabSplit) > 0 {
				softwareName := tabSplit[0]
				softwareVersion := tabSplit[1]
				softwareVendor := tabSplit[2]
				installStruct := Software{softwareName, softwareVersion, softwareVendor}
				softwareList = append(softwareList, installStruct)
			}
		}
	}*/

	return softwareList
}
