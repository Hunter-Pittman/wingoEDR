package honeymonitor

import (
	"fmt"
	"strings"
	"time"
	//"wingoEDR/common"

	//"wingoEDR/logger"

	"io/ioutil"
	"os"

	"github.com/djherbis/times"
	"go.uber.org/zap"
	"encoding/hex"
	"crypto/sha1"
)

var (
	modTimeDirectory = os.Getenv("TEMP") + "\\HoneyDir\\" // Directory where modtime files will be stored
	modTimeOrigFile = modTimeDirectory + "modtime.orig"
	modTimeNewFile = modTimeDirectory + "modtime.new"
	
)

type FileAttribs struct {
	filepath   string
	modTime    time.Time
	accessTime time.Time
}

// New & Improved Func
// DO NOT add / to end of path (Acceptable path example: /home/user/honeydir)
func enumerateFiles(directoriesToMonitor []string) []string {
	// iterate through each target directory, determine whether or not the directory is actually a directory or not
	var targetMonitorPaths []string
	for _, filepath := range directoriesToMonitor {
		handle, err := os.Stat(filepath)
		if err != nil || os.IsNotExist(err) {
			zap.S().Fatal("Specified filepath when running enumerateFiles() function is nonexistent: " + err.Error() + filepath)
		}
		// determine whether or not path is a directory or not 
		if handle.IsDir() {
			fileList, err := ioutil.ReadDir(filepath)
			if err != nil {
				zap.S().Error("Cannot read directory specified: " + filepath + err.Error())
			}
			for _, file := range fileList {
				targetMonitorPaths = append(targetMonitorPaths, filepath + "/" + file.Name())
			}
		} else {
			// error not showing up on terminal
			zap.S().Error("Specified filepath when running enumerateFiles() isn't an actual directory: " + filepath)
		}
	}
	return targetMonitorPaths
}


// create & return fileAttribs structure for 1 file
func getFileAttribs(filePath string) FileAttribs {
	data, err := times.Stat(filePath)
	if err != nil {
		if strings.Contains(err.Error(), "cannot find the file") {
			// sign that the file was deleted
			zap.S().Error("1 honey file was likely deleted! Sending alert:" + err.Error())
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
func getTimes(fileList []string) []FileAttribs  {
	files := []FileAttribs{}
	for _, file := range fileList {
		fileAttributes := getFileAttribs(file)
		files = append(files, fileAttributes)
	}
	return files
}

// Check if directory already exists. If it does, check if service file exists. If not, create it. Otherwise create new service file & compare against old one
// Todo: implement compression
// Return both filepaths
func CreateDirMonitor(directories []string) {
	_, err := os.Stat(modTimeOrigFile)
	if os.IsNotExist(err) {
		zap.S().Info(modTimeOrigFile + " does not exist! Creating...")
		_ = os.Mkdir(modTimeDirectory, 0755) //may want to change numbers to restrict permissions only to owner of directory to prevent modification (EX: 0600)
		fileList := enumerateFiles(directories)
		origTimes := getTimes(fileList)
		fileHandle, err := os.Create(modTimeOrigFile)
		if err != nil {
			zap.S().Fatal("Error creating the original modification time file: + " + err.Error())
		}
		for _, line := range origTimes {
			strToWrite := line.filepath + " " + line.modTime.String() + " " + line.accessTime.String()
			fileHandle.Write([]byte(strToWrite + "\n"))
		}
		fileHandle.Close()
	} else {
		// create new file with new times
		zap.S().Info(modTimeOrigFile + " exists! Creating new times file to compare to old...")
		fileList := enumerateFiles(directories)
		newTimes := getTimes(fileList)
		fileHandle, err := os.Create(modTimeNewFile)
		if err != nil {
			zap.S().Error("Could not create the new modification time file: " + err.Error()
		}
		for _, line := range newTimes {
			strToWrite := line.filepath + " " + line.modTime.String() + " " + line.accessTime.String()
			fileHandle.Write([]byte(strToWrite + "\n"))
		}
		fileHandle.Close()
		
		// compare old file with new
		// send alert if function returns false
		if compareFiles(modTimeOrigFile, modTimeNewFile) {
			zap.S().Info("Honeypot files untouched.")
		} else {
			zap.S().Fatal("Honeypot files touched! Potential intrusion!")
	}
	}
}

// Used to compare two files using sha1 hashes
func compareFiles(filepath1 string, filepath2 string) bool {
	// ensure both files exist
	_, err1 := os.Stat(filepath1)
	_, err2 := os.Stat(filepath2)
	if os.IsNotExist(err1) || os.IsNotExist(err2) {
		zap.S().Fatal("One of the two specified filepaths passed to compareFiles() doesn't exist: " + filepath1 + " or " + filepath2)
	}
	// read file data
	filepath1Content, err := ioutil.ReadFile(filepath1)
	if err != nil {
		zap.S().Fatal("Cannot read: " + filepath1)
	}
	filepath2Content, err := ioutil.ReadFile(filepath2)
	if err != nil {
		zap.S().Fatal("Cannot read: " + filepath2)
	}
	
	// hash & compare file data
	filepath1Byte := sha1.Sum(filepath1Content)
	filepath1Hash := hex.EncodeToString(filepath1Byte[:])
	
	filepath2Byte := sha1.Sum(filepath2Content)
	filepath2Hash := hex.EncodeToString(filepath2Byte[:])

	if filepath2Hash != filepath1Hash {
		return false
	} else {
		return true
	}
}
