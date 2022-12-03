package firewall

import (
	wapi "github.com/iamacarpet/go-win64api"
	"go.uber.org/zap"
)

type FirewallList struct {
	Name            string
	Description     string
	ApplicationName string
	ServiceName     string
	LocalPorts      string
	RemotePorts     string
	LocalAddresses  string
	RemoteAddresses string
	Profile         int32
}

func FirewallLister() []FirewallList {
	fireslice := make([]FirewallList, 0)
	svc, err := wapi.FirewallRulesGet()
	if err != nil {
		zap.S().Error("Getting services failed!")
	}
	for _, v := range svc {

		helium := FirewallList{
			Name:            v.Name,
			Description:     v.Description,
			ApplicationName: v.ApplicationName,
			ServiceName:     v.ServiceName,
			LocalPorts:      v.LocalPorts,
			RemotePorts:     v.RemotePorts,
			LocalAddresses:  v.LocalAddresses,
			RemoteAddresses: v.RemoteAddresses,
			Profile:         v.Profiles}

		fireslice = append(fireslice, helium)

	}
	return fireslice

}

func FWruleadd(apprulename string, appPath string, portString string, profile32 int32) bool {

	added, err := wapi.FirewallRuleAddApplication(
		apprulename,
		"App Rule Long Description.",
		appPath,
		portString,
		profile32,
	)

	if err != nil {
		zap.S().Error("Getting services failed!")
	}

	return added

}

func FWruleremove(remrulename string) bool {
	removed, err := wapi.FirewallRuleDelete(remrulename)

	if err != nil {
		zap.S().Error("Getting services failed!")
	}

	return removed
}

func FWDisable(profile32 int32) bool {
	disabled, err := wapi.FirewallDisable(profile32)

	if err != nil {
		zap.S().Error("Getting services failed!")
	}

	return disabled
}

func FWEnable(profile32 int32) bool {
	enabled, err := wapi.FirewallEnable(profile32)

	if err != nil {
		zap.S().Error("Getting services failed!")
	}

	return enabled
}
