package shares

import (
	"os/exec"
	"strings"

	"github.com/StackExchange/wmi"
	"go.uber.org/zap"
)

type Win32_Share struct {
	Name string
	Path string
}

type ShareAttributes struct {
	Name        string
	Path        string
	Permissions string
}

var (
	logger, _ = zap.NewProduction()
)

// return string with permissions for specified share path
func getSharePermissions(sharePath string) string {
	// if sharePath arg is empty, we return empty string
	if len(sharePath) != 0 {

		toExecute := "Get-Acl -Path '" + sharePath + "' | Select-Object -ExpandProperty Access | Select-Object -Property FileSystemRights, AccessControlType, Identityreference | ConvertTo-CSV"
		unformattedPerms, err := exec.Command("powershell.exe", toExecute).Output()
		if err != nil {
			logger.Info(err.Error())
			//fmt.Println(err)
		}
		//replacing ints with respective permissions
		formattedPerms1 := strings.Replace(string(unformattedPerms), "\"268435456\"", "\"FullControl\"", -1)
		formattedPerms2 := strings.Replace(formattedPerms1, "\"-536805376\"", "\"Modify,Synchronize\"", -1)
		finalFormattedPerms := strings.Replace(formattedPerms2, "\"-1610612736\"", "\"ReadAndExecute\"", -1)
		return finalFormattedPerms

	} else {
		return ""
	}
}

// return slice/list of share names
func ListSharesWMI() []ShareAttributes {
	var dst []Win32_Share
	shares := []ShareAttributes{}
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)
	if err != nil {
		logger.Error(err.Error())
		//log.Fatal(err)
	}
	// can get each name of a share & query for its permissions
	for _, shareAttrib := range dst {
		permissions := getSharePermissions(shareAttrib.Path)
		var share = ShareAttributes{shareAttrib.Name, shareAttrib.Path, permissions}
		shares = append(shares, share)

	}
	return shares
}
