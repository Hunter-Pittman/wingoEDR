package softwaremanager

import (
    "fmt"
    wapi "github.com/iamacarpet/go-win64api"
)

type Software struct {
	Name	string
	Version	string
}

func GetSoftware() []Software {
    sw, err := wapi.InstalledSoftwareList()
    if err != nil {
        common.ErrorHandler(err)
	}

	softwareList := []Software{}
	for _, s := range sw {
		softwareInfo := Software{s.Name(), s.Version()}
		softwareList = append(softwareList, softwareInfo)
	}
	return softwareList
}

