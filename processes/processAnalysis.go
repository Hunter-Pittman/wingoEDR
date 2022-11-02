package processes

import (
	"fmt"
	"regexp"
)

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

// Handle network analysis from seperate module
