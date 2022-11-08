package systemhealth

import (
	//"fmt"
	"time"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/uptime"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/cpu"
	//"log"
	"go.uber.org/zap"
	"math"
)

var (
	logger, _ = zap.NewProduction()
)

type MemInfo struct {
	Total	float64
	Used	float64
	Avail	float64
}

type DiskData struct {
	Total	float64
	Used	float64
	UsedPercent float64
	Free	float64
	
}

type SystemStats struct {
	MemoryInfo	MemInfo
	DiskInfo	DiskData
	NetworkInfo	net.IOCountersStat
	Uptime		time.Duration
	CpuUsage	[]float64
}

func getMemoryStats() MemInfo {
	memStats, err := memory.Get()	
	if err != nil {
		//log.Fatal(err)
		logger.Error(err.Error())
	}
	totalFloat64 := ConvertBytesToGigabytes(memStats.Total)
	usedFloat64 := ConvertBytesToGigabytes(memStats.Used)
	freeFloat64 := ConvertBytesToGigabytes(memStats.Free)

	return MemInfo{totalFloat64, usedFloat64, freeFloat64}
}


func getSystemUptime() time.Duration {
	uptime, err := uptime.Get()
	if err != nil {
		//log.Fatal(err)
		logger.Error(err.Error())
	}
	return uptime
}

/*returns network data
	BytesSent
	BytesRecv
	PacketsSent
	PacketsRecv
https://pkg.go.dev/github.com/shirou/gopsutil/v3@v3.22.9/net#IOCountersStat
*/
func getNetworkData() net.IOCountersStat {
	networkData, err := net.IOCounters(false) //Change to true if you want to separate network data for each interface
	if err != nil {
		logger.Info(err.Error())
		//fmt.Println(err)
	}
	return networkData[0]
}

func getDiskUsageData() DiskData {
	out, err := disk.Usage("/")
	if err != nil {
		logger.Error(err.Error())
		//log.Fatal(err)
	}
	return DiskData{ConvertBytesToGigabytes(out.Total), ConvertBytesToGigabytes(out.Used), out.UsedPercent, ConvertBytesToGigabytes(out.Free)}
}

func ConvertBytesToGigabytes(bytes uint64) float64 {
	gigabytes := float64(bytes)/1073741824
	return math.Round(gigabytes)
}

//may not work, requires more testing
func getCpuUsage() []float64 {
	out, err := cpu.Percent(0, false)
	if err != nil {
		logger.Error(err.Error())
		//log.Fatal(err)
	}
	return out
}

// compiles all functions & returns struct with all necessary data
// refer to SystemStats for what it returns
func GetSystemHealth() SystemStats {
	memoryStats := getMemoryStats()
	uptime := getSystemUptime()
	networkStats := getNetworkData()
	diskUsageStats := getDiskUsageData()
	cpuUsageStats := getCpuUsage()
	return SystemStats{memoryStats, diskUsageStats, networkStats, uptime, cpuUsageStats}
}
