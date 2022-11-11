package honeytoken

import (
	//"fmt"
	"fmt"
	"time"
	"wingoEDR/common"

	//"log"
	"encoding/base64"
	"io/ioutil"

	"github.com/djherbis/times"
	"go.uber.org/zap"
)

type FileAttribs struct {
	Filename   string
	ModTime    time.Time
	AccessTime time.Time
	FileHash   string
	FileData   string
}

// create & return fileAttribs structure for 1 file
func GetFileAttribs(filepath string) FileAttribs {
	data, err := times.Stat(filepath)
	if err != nil {
		//log.Fatal(err)
		zap.S().Error(err.Error())
	} else {
		// create file struct
		fileData, _ := ioutil.ReadFile(filepath)
		fileDataB64 := base64.StdEncoding.EncodeToString([]byte(fileData))
		fileHash := common.GenerateSha1Hash(string(fileData))
		var honeyFileAttribs = FileAttribs{filepath, data.ModTime(), data.AccessTime(), fileHash, fileDataB64}
		return honeyFileAttribs
	}
	var honeyFileAttribs = FileAttribs{"", time.Now(), time.Now(), "", ""}
	return honeyFileAttribs
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
			incident := common.Incident{
				Name:     fmt.Sprintf("Honeytoken access violation on: %v", filepath),
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
			//fmt.Println("[+] File untouched, resuming sleep...")
			zap.S().Info("[+] File untouched, resuming sleep...")
		}
	}
}
