package shares

import (
	"os/exec"
	"log"
	"fmt"
	"github.com/StackExchange/wmi"
	"strings"
)

type Win32_Share struct {
	Name string
	Path string
}

type shareAttributes struct {
	Name	string
	Path	string
	Permissions string
}



//return string with permissions for specified share path
func getSharePermissions(sharePath string) string {
	// if sharePath arg is empty, we return empty string
	if len(sharePath) != 0 {
	
		toExecute := "Get-Acl -Path '" + sharePath + "' | Select-Object -ExpandProperty Access | Select-Object -Property FileSystemRights, AccessControlType, Identityreference | ConvertTo-CSV"
		unformattedPerms, err := exec.Command("powershell.exe", toExecute).Output()
		if err != nil {
			fmt.Println(err)
		}
		//replacing ints with respective permissions
		formattedPerms1 := strings.Replace(string(unformattedPerms), "\"268435456\"", "\"FullControl\"",-1)
		formattedPerms2 := strings.Replace(formattedPerms1, "\"-536805376\"", "\"Modify,Synchronize\"",-1)
		finalFormattedPerms := strings.Replace(formattedPerms2, "\"-1610612736\"", "\"ReadAndExecute\"",-1)
		return finalFormattedPerms

	} else {
		return ""
	}
}



//return slice/list of share names
func listSharesWMI() []shareAttributes {
	var dst []Win32_Share
	shares := []shareAttributes{}
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)
	if err != nil {
		log.Fatal(err)
	}
	// can get each name of a share & query for its permissions
	for _, shareAttrib := range dst {
		permissions := getSharePermissions(shareAttrib.Path)
		var share = shareAttributes{shareAttrib.Name, shareAttrib.Path, permissions}
		shares = append(shares, share)
		
	}
	return shares
}

