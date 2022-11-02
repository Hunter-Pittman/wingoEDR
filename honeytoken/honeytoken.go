package honeytoken

import (
	//"fmt"
	"time"
	//"log"
	"go.uber.org/zap"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"github.com/djherbis/times"
	"encoding/base64"
)

type FileAttribs struct {
	Filename	string
	ModTime		time.Time
	AccessTime	time.Time
	FileHash	string
	FileData	string
}

var (
	logger, _ = zap.NewProduction()
)

//create & return fileAttribs structure for 1 file
func GetFileAttribs(filepath string) FileAttribs {
	data, err := times.Stat(filepath)
	if err != nil {
		//log.Fatal(err)
		logger.Error(err.Error())
	} else {
		// create file struct
		fileData, _ := ioutil.ReadFile(filepath)
		fileDataB64 := base64.StdEncoding.EncodeToString([]byte(fileData))
		fileHash := GenerateSha1Hash(string(fileData))
		var honeyFileAttribs = FileAttribs{filepath, data.ModTime(), data.AccessTime(), fileHash, fileDataB64}
		return honeyFileAttribs
	}
	var honeyFileAttribs = FileAttribs{"", time.Now(), time.Now(), "", ""}
	return honeyFileAttribs
}

func GenerateSha1Hash(data string) string {
	dataHashByte := sha1.Sum([]byte(data))
	dataHashStr := hex.EncodeToString(dataHashByte[:])
	return dataHashStr
}

// Get file attributes (modification time, access time, filename), store in honeyAttribs1
// Sleep for x amount of time
// Get file attributes & store it in honeyAttrib2
// Compare honeyAttrib1 & honeyAttrib2, then determine whether or not file tampering has occurred
func MonitorHoneyFile(filepath string) {
	
	for {
		honeyAttribs1 := GetFileAttribs(filepath)
		time.Sleep(2 * time.Second)
		honeyAttribs2 := GetFileAttribs(filepath)
		if honeyAttribs1.ModTime != honeyAttribs2.ModTime || honeyAttribs1.AccessTime != honeyAttribs2.AccessTime {
			logger.Warn("[!] File has been accessed and/or modified! Sending alert!")
			//log.Fatal("[!] File has been accessed and/or modified! Sending alert!")
		} else {
			//fmt.Println("[+] File untouched, resuming sleep...")
			logger.Info("[+] File untouched, resuming sleep...")
		}
	}
}

