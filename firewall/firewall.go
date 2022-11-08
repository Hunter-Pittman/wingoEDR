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
}

func Firewalllister() []FirewallList {
	fireslice := make([]FirewallList, 0)
	svc, err := wapi.FirewallRulesGet()
	if err != nil {
		zap.S().Error("Getting services failed!")
	}
	for _, v := range svc {

		hydrogen := FirewallList{
			name:        v.Name,
			desc:        v.Description,
			accname:     v.ApplicationName,
			servname:    v.ServiceName,
			localports:  v.LocalPorts,
			remoteports: v.RemotePorts,
			localaddr:   v.LocalAddresses,
			remoteaddr:  v.RemoteAddresses}

		fireslice = append(fireslice, hydrogen)

	}
	return fireslice

}
