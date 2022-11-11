package honeydirectory

import (
	"fmt"
	"strings"
	"time"
	"wingoEDR/common"

	//"wingoEDR/logger"

	"io/ioutil"
	"os"

	"github.com/djherbis/times"
	"go.uber.org/zap"
)

type FileAttribs struct {
	filename   string
	modTime    time.Time
	accessTime time.Time
}

// path is directory path without appended "/" (EX: /home/user/test)
func enumerateFiles(path string) []string {
	// verify file/dir exists
	fileHandle, err := os.Stat(path)
	if err != nil {
		//path doesn't exist
		zap.S().Error("Specified directory path when enumerating files is nonexistent: " + err.Error())
		//logger.S()("Specified directory path when enumerating files is nonexistent: ", err)
	}
	// if it's dir, enumerate it & return all files
	if fileHandle.IsDir() {
		fileList, err := ioutil.ReadDir(path)
		fileListSlice := []string{}
		if err != nil {
			zap.S().Error("Cannot read directory specified:" + err.Error())
			//zap.S().Fatal("Cannot read directory specified: ", err)
		}
		for _, file := range fileList {
			fileListSlice = append(fileListSlice, path+"/"+file.Name())
		}
		return fileListSlice
	} else {
		// return slice with only single file path
		filepath := []string{path}
		return filepath
	}
}

// create & return fileAttribs structure for 1 file
func getFileAttribs(filePath string) FileAttribs {
	data, err := times.Stat(filePath)
	if err != nil {
		if strings.Contains(err.Error(), "cannot find the file") {
			// sign that the file was deleted
			zap.S().Error("1 honey file was likely deleted! Sending alert:" + err.Error())
			//zap.S().Warn("1 honey file was likely deleted! Sending alert:", err)
		}
	} else {
		// create file struct
		var honeyDirAttribs = FileAttribs{filePath, data.ModTime(), data.AccessTime()}
		return honeyDirAttribs
	}
	var honeyFileAttribs = FileAttribs{"", time.Now(), time.Now()}
	return honeyFileAttribs
}

// compile list & return list of all file access times
// then sleep for x time, compile another list & we loop through the list & compare each file access time
func getTimes(fileList []string) []FileAttribs {
	files := []FileAttribs{}
	for _, file := range fileList {
		fileAttributes := getFileAttribs(file)
		files = append(files, fileAttributes)
	}
	return files
}

// Example: monitorDirectory("C:\\Users\\User\\Desktop\\honeytoken", 2) to monitor honeytoken & sleep for 2 seconds
func MonitorHoneyDirectory(directory string, secondsToSleep time.Duration) {
	fileList := enumerateFiles(directory)
	origTimes := getTimes(fileList)
	for {
		time.Sleep(time.Second * secondsToSleep)
		newTimes := getTimes(fileList)
		// new file added or file deleted
		if len(origTimes) != len(newTimes) {
			zap.S().Warn("File added or deleted! Sending alert!")
		} else {
			for index, file := range origTimes {
				if file.modTime != newTimes[index].modTime || file.accessTime != newTimes[index].accessTime {
					zap.S().Warn("Honey file accessed/modified! Sending alert!")

					incident := common.Incident{
						Name:     fmt.Sprintf("Honey directory access violation on: %v", directory),
						User:     "Bob",
						Process:  "",
						RemoteIP: "",
						Cmd:      "",
					}

					alert := common.Alert{
						Host:     common.GetSerialScripterHostName(),
						Incident: incident,
					}
					common.IncidentAlert(alert)

				} else {
					zap.S().Info("[+] Files untouched. Sleeping...")
				}
			}
		}
	}
}
