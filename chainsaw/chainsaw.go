package chainsaw

import (
	"os"
	"os/exec"
	"wingoEDR/common"
)

func ScanAll() {
	//cmdOutput, err := exec.Command("C:\\Users\\Hunter Pittman\\Documents\\repos\\wingoEDR\\externalresources\\chainsaw\\chainsaw_x86_64-pc-windows-msvc.exe").Output()
	cmdOutput, err := exec.Command("C:\\Users\\Hunter Pittman\\Documents\\repos\\wingoEDR\\externalresources\\chainsaw\\chainsaw_x86_64-pc-windows-msvc.exe", "--no-banner", "hunt", "C:\\Windows\\System32\\winevt\\Logs\\", "-s", "C:\\Users\\Hunter Pittman\\Documents\\repos\\wingoEDR\\externalresources\\chainsaw\\sigma\\rules\\windows", "--mapping", "C:\\Users\\Hunter Pittman\\Documents\\repos\\wingoEDR\\externalresources\\chainsaw\\mappings\\sigma-event-logs-all.yml", "--json").Output()
	common.ErrorHandler(err)

	f, err := os.Create("C:\\test1.txt")
	common.ErrorHandler(err)
	defer f.Close()

	_, err2 := f.WriteString(string(cmdOutput))
	common.ErrorHandler(err2)

	println(string(cmdOutput))
}
