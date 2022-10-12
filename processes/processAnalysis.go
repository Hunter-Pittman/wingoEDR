package processes

import (
	"log"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessInfo struct {
	name string
	ppid *process.Process
	pid  int32
	exe  []string
	cwd  string
	user string
}

func GetAllProcesses() []ProcessInfo {

	processList, err := process.Processes()
	if err != nil {
		log.Println("Getting processes failed!")
	}
	result := make([]ProcessInfo, len(processList))
	for i := range processList {

		ppid, err := processList[i].Parent()
		if err != nil {
			log.Println("Error getting parent process: ", err)
		}

		exe, err := processList[i].CmdlineSlice()
		if err != nil {
			log.Println("Error getting command line arguments: ", err)
		}

		cwd, err := processList[i].Cwd()
		if err != nil {
			log.Println("Error getting process current working directory: ", err)
		}

		name, err := processList[i].Name()
		if err != nil {
			log.Println("Error getting process name: ", err)
		}

		user, err := processList[i].Username()
		if err != nil {
			log.Println("Error getting username: ", err)
		}

		result[i] = ProcessInfo{
			name: name,
			ppid: ppid,
			pid:  processList[i].Pid,
			exe:  exe,
			cwd:  cwd,
			user: user,
		}
	}

	return result
}
