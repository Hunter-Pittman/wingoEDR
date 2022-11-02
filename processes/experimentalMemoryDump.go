package processes

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sys/windows"
)

// Don't use this file still working on code

// ProcessInformation wrap basic process information and memory dump in a structure
type ProcessInformation struct {
	PID         uint32
	ProcessName string
	ProcessPath string
	MemoryDump  []byte
}

/*
func GetInfo() {
	var procID uint32 = 16564
	procHandle, err := getProcessHandle(procID, windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ)
	if err != nil {
		fmt.Printf("%v", err)
	}

	procFilename, modules, err := getProcessModulesHandles(procHandle)
	if err != nil {
		fmt.Printf("%v", err)
	}

	dumpModuleMemory(procHandle, modules)
}

*/

// GetProcessMemory return a process memory dump based on its handle
func getProcessMemory(pid uint32, handle windows.Handle, verbose bool) (ProcessInformation, []byte, error) {

	procFilename, modules, err := getProcessModulesHandles(handle)
	if err != nil {
		return ProcessInformation{}, nil, fmt.Errorf("Unable to get PID %d memory: %s", pid, err.Error())
	}

	for _, moduleHandle := range modules {
		if moduleHandle != 0 {
			moduleRawName, err := GetModuleFileNameEx(handle, moduleHandle, 512)
			if err != nil {
				return ProcessInformation{}, nil, err
			}
			moduleRawName = bytes.Trim(moduleRawName, "\x00")
			modulePath := strings.Split(string(moduleRawName), "\\")
			moduleFileName := modulePath[len(modulePath)-1]

			if procFilename == moduleFileName {
				return ProcessInformation{PID: pid, ProcessName: procFilename, ProcessPath: string(moduleRawName)}, dumpModuleMemory(handle, moduleHandle, verbose), nil
			}
		}
	}

	return ProcessInformation{}, nil, fmt.Errorf("Unable to get PID %d memory: no module corresponding to process name", pid)
}

// GetProcessesList return PID from running processes
func getProcessesList() (procsIds []uint32, bytesReturned uint32, err error) {
	procsIds = make([]uint32, 2048)
	err = windows.EnumProcesses(procsIds, &bytesReturned)
	return procsIds, bytesReturned, err
}

// GetProcessHandle return the process handle from the specified PID
func getProcessHandle(pid uint32, desiredAccess uint32) (handle windows.Handle, err error) {
	handle, err = windows.OpenProcess(desiredAccess, false, pid)
	return handle, err
}

// getProcessModulesHandles list modules handles from a process handle
func getProcessModulesHandles(procHandle windows.Handle) (processFilename string, modules []syscall.Handle, err error) {
	var processRawName []byte
	processRawName, err = GetProcessImageFileName(procHandle, 512)
	if err != nil {
		return "", nil, err
	}
	processRawName = bytes.Trim(processRawName, "\x00")
	processPath := strings.Split(string(processRawName), "\\")
	processFilename = processPath[len(processPath)-1]

	modules, err = EnumProcessModules(procHandle, 32)
	if err != nil {
		return "", nil, err
	}

	return processFilename, modules, nil
}

// dumpModuleMemory dump a process module memory and return it as a byte slice
func dumpModuleMemory(procHandle windows.Handle, modHandle syscall.Handle, verbose bool) []byte {
	moduleInfos, err := GetModuleInformation(procHandle, modHandle)
	if err != nil && verbose {
		zap.S().Error("Dumping memory modules failed: ", err)
	}

	memdump, err := ReadProcessMemory(procHandle, moduleInfos.BaseOfDll, uintptr(moduleInfos.SizeOfImage))
	if err != nil && verbose {
		zap.S().Error("Reading memory modules failed: ", err)
	}

	memdump = bytes.Trim(memdump, "\x00")
	return memdump
}

// WriteProcessMemoryToFile try to write a byte slice to the specified directory
func writeProcessMemoryToFile(path string, file string, data []byte) (err error) {
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0600); err != nil {
			return err
		}
	}

	if err := os.WriteFile(path+"/"+file, data, 0644); err != nil {
		return err
	}

	return nil
}
