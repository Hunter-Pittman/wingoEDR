package softwaremanager

import (
	wapi "github.com/iamacarpet/go-win64api"
	"go.uber.org/zap"
)

type Software struct {
	Name    string
	Version string
}

func GetSoftware() []Software {
	sw, err := wapi.InstalledSoftwareList()
	if err != nil {
		zap.S().Error("Error getting software list: ", err)
	}

	softwareList := []Software{}
	for _, s := range sw {
		softwareInfo := Software{s.Name(), s.Version()}
		softwareList = append(softwareList, softwareInfo)
	}
	return softwareList
}
