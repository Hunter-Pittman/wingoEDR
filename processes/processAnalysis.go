package processes

import (
	"fmt"
	"regexp"
	"strings"
)

type SuspiciousCmd struct {
	RanCommand     string
	MatchedKeyword string
}

var processes, _ = GetAllProcesses()

func charateristicsAnalysis(processList *[]ProcessInfo) {

}

func MaliciousKnownPrograms() {
	maliciousProgramPattern := regexp.MustCompile(`bash|sh|.php$|base64|nc|ncat|shell|^python|telnet|ruby`)

	for i := range processes {
		if maliciousProgramPattern.MatchString(processes[i].Name) {
			fmt.Printf("user running a shell %v\n", processes[i].Name)
		}
	}

}

func WindowsFindDeviousCmdParams(cmd string) SuspiciousCmd {
	suspiciousCLParams := []string{
		//https://redcanary.com/threat-detection-report/techniques/windows-command-shell/
		//https://arxiv.org/pdf/1804.04177v2.pdf
		"^",
		"=",
		"%",
		"!",
		"[",
		"(",
		";",
		"http",
		"https",
		"echo",
		"cmd.exe /c",
		"cmd /c",
		"write-host",
		"bypass",
		"exec",
		"create",
		"dumpcreds",
		"downloadstring",
		"invoke-command",
		"getstring",
		"webclient",
		"nop",
		"hidden",
		"encodedcommand",
		"-en",
		"-enc",
		"-enco",
		"downloadfile",
		"iex",
		"replace",
		"wscript.shell",
		"windowstyle",
		"comobject",
		"reg",
		"autorun",
		"psexec",
		"lsadump",
		"wmic",
		"schtask",
		"net",
		"fsutil",
		"dsquery",
		"netsh",
		"del",
		"taskkill",
		"uploadfile",
		"invoke-wmi",
		"enumnetworkdrives",
		"procdump",
		"get-wmiobject",
		"sc",
		"cimv2",
		"-c",
		"certutil",
		"new-itemproperty",
		"invoke-expression",
		"invoke-obfuscation",
		"nop",
		"invoke-webrequest",
		"reflection",
	}

	if len(cmd) != 0 {
		for _, knownParam := range suspiciousCLParams {
			lowerCaseKnownParam := strings.ToLower(knownParam)
			lowerCaseCmd := strings.ToLower(cmd)
			if strings.Contains(lowerCaseCmd, lowerCaseKnownParam) {
				//fmt.Println("[+] Potentially malicious command found:" + cmd)
				//fmt.Println("[+] Keyword match:" + lowerCaseKnownParam)
				return SuspiciousCmd{cmd, lowerCaseKnownParam}
			}
		}
	}
	return SuspiciousCmd{"", ""}
}

// Handle network analysis from seperate module
