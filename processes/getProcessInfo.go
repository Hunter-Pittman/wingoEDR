package processes

import (
	"errors"

	"github.com/shirou/gopsutil/v3/process"
	"go.uber.org/zap"
)

// Additoanl data functions needed:
// 1. Get name via PID
// 2. Get PID via name

type ProcessInfo struct {
	Name               string `default:"no name"`
	Ppid               *process.Process
	Pid                int32
	Exe                []string
	Cwd                string
	User               string
	NetworkConnections []ConnectionInfo
}

type ConnectionInfo struct {
	NetType       uint32
	LocalAddress  string
	LocalPort     uint32
	RemoteAddress string
	RemotePort    uint32
	Status        string
}

func GetAllProcesses() []ProcessInfo {

	processList, err := process.Processes()
	if err != nil {
		zap.S().Error("Getting processes failed!")
	}
	result := make([]ProcessInfo, len(processList))
	for i := range processList {

		ppid, err := processList[i].Parent()
		if err != nil {
			if err.Error() == "process does not exist" {
			} else {
				zap.S().Warn("Error getting parent process: ", err)
			}
		}

		exe, err := processList[i].CmdlineSlice()
		if err != nil {
			zap.S().Warn("Error getting command line arguments: ", err)
		}

		cwd, err := processList[i].Cwd()
		if err != nil {
			zap.S().Warn("Error getting process current working directory: ", err)
		}

		name, err := processList[i].Name()
		if err != nil {
			zap.S().Warn("Error getting process name: ", err)
		}

		user, err := processList[i].Username()
		if err != nil {
			zap.S().Warn("Error getting username: ", err)
		}

		connections, err := GetProcessConnections(processList[i])
		if err != nil {
			if err.Error() == "There are no connections for this process" {
			} else {
				zap.S().Warn("Error getting process netowrk connections: ", err)
			}
		}

		result[i] = ProcessInfo{
			Name:               name,
			Ppid:               ppid,
			Pid:                processList[i].Pid,
			Exe:                exe,
			Cwd:                cwd,
			User:               user,
			NetworkConnections: connections,
		}
	}

	return result
}

func GetProcessConnections(process *process.Process) ([]ConnectionInfo, error) {
	connections, err := process.Connections()
	var noConnections []ConnectionInfo
	if err != nil {
		return noConnections, errors.New("Error getting connections for this process")
	}
	numberOfConnections := len(connections)
	currentConnections := make([]ConnectionInfo, numberOfConnections)

	if numberOfConnections == 0 {
		return noConnections, errors.New("There are no connections for this process")
	}

	for i := range connections {

		currentConnections[i] = ConnectionInfo{
			NetType:       connections[i].Type,
			LocalAddress:  connections[i].Laddr.IP,
			LocalPort:     connections[i].Laddr.Port,
			RemoteAddress: connections[i].Raddr.IP,
			RemotePort:    connections[i].Raddr.Port,
			Status:        connections[i].Status,
		}
	}

	return currentConnections, nil
}
