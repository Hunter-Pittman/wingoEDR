package processes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/schollz/closestmatch"
)

type SuspiciousCmd struct {
	RanCommand     string
	MatchedKeyword string
	ProcessName    string
	Pid            int32
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

func FindDeviousCmdParams() []SuspiciousCmd {
	suspiciousCLParams := []string{
		//https://redcanary.com/threat-detection-report/techniques/windows-command-shell/
		//https://arxiv.org/pdf/1804.04177v2.pdf
		"^",
		"=",
		"%",
		"!",
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
		"nc",
		"ncat",
		"python",
		"telnet",
	}

	bagSizes := []int{2}
	cm := closestmatch.New(suspiciousCLParams, bagSizes)

	badProcesses := make([]SuspiciousCmd, 0)

	for i := range processes {
		if len(processes[i].Exe) != 0 {
			for _, knownParam := range suspiciousCLParams {
				lowerCaseKnownParam := strings.ToLower(knownParam)
				for x := range processes[i].Exe {
					lowerCaseCmd := strings.ToLower(processes[i].Exe[x])
					if cm.Closest(lowerCaseCmd) != "" {
						accuracy := cm.AccuracyMutatingWords()
						fmt.Printf("Accuracy stuff %v", accuracy)
						fullCommand := processes[i].Name
						for y := range processes[i].Exe {
							fullCommand += " " + processes[i].Exe[y]
						}

						thing := SuspiciousCmd{fullCommand, lowerCaseKnownParam, processes[i].Name, processes[i].Pid}
						badProcesses = append(badProcesses, thing)
					}
					//if strings.Contains(lowerCaseCmd, lowerCaseKnownParam) {}
				}

			}
		}
	}
	return badProcesses
}

// Handle network analysis from seperate module

func DeviousCommandPartialSearchTest() []SuspiciousCmd {
	suspiciousCLParams := []string{
		//https://redcanary.com/threat-detection-report/techniques/windows-command-shell/
		//https://arxiv.org/pdf/1804.04177v2.pdf
		"^",
		"=",
		"%",
		"!",
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
		//		"windowstyle",
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
		"nc.exe",
		"ncat",
		"python",
		"telnet",
		"foxit",
	}
	bagSizes := []int{2}
	cm := closestmatch.New(suspiciousCLParams, bagSizes)

	badProcesses := make([]SuspiciousCmd, 0)

	for i := range processes {
		closestHolder := cm.Closest(processes[i].Name)
		if closestHolder != "" {
			fullCommand := processes[i].Name
			for y := range processes[i].Exe {
				fullCommand += " " + processes[i].Exe[y]
			}
			accuracy := cm.AccuracyMutatingWords()
			if accuracy > 75.0 {
				fmt.Printf("Name Match found. Keyword: %v Command: %v \n", closestHolder, fullCommand)
				fmt.Printf("Accuracy Rating: %v \n", cm.AccuracyMutatingWords())
			}

			closestHolder = ""
		}

		for x := range processes[i].Exe {
			closestHolder := cm.Closest(processes[i].Exe[x])
			if closestHolder != "" {
				fullCommand := processes[i].Name
				for y := range processes[i].Exe {
					fullCommand += " " + processes[i].Exe[y]
				}
				accuracy := cm.AccuracyMutatingWords()
				if accuracy > 80.0 {
					fmt.Printf("Match found. Keyword: %v Command: %v \n", closestHolder, fullCommand)
					fmt.Printf("Accuracy Rating: %v \n", cm.AccuracyMutatingWords())
				}

				closestHolder = ""
			}
		}
	}

	return badProcesses
}
