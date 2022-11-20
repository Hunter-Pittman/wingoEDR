package firewall

import (
	wapi "github.com/iamacarpet/go-win64api"
	"go.uber.org/zap"
)

type FirewallList struct {
	name        string
	desc        string
	accname     string
	servname    string
	localports  string
	remoteports string
	localaddr   string
	remoteaddr  string
	profile     int32
}

func Firewalllister() []FirewallList {
	fireslice := make([]FirewallList, 0)
	svc, err := wapi.FirewallRulesGet()
	if err != nil {
		zap.S().Error("Getting services failed!")
	}
	for _, v := range svc {

		helium := FirewallList{
			name:        v.Name,
			desc:        v.Description,
			accname:     v.ApplicationName,
			servname:    v.ServiceName,
			localports:  v.LocalPorts,
			remoteports: v.RemotePorts,
			localaddr:   v.LocalAddresses,
			remoteaddr:  v.RemoteAddresses,
			profile:     v.Profiles}

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
