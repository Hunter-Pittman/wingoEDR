package honeytoken

import (
	"fmt"
	"time"
	"log"
	"crypto/sha1"
	"encoding/hex"
	"github.com/djherbis/times"
)

type fileAttribs struct {
	filename	string
	modTime		time.Time
	accessTime	time.Time
}

//create & return fileAttribs structure for 1 file
func getFileAttribs(filepath string) fileAttribs {
	data, err := times.Stat(filepath)
	if err != nil {
		log.Fatal(err)
	} else {
		// create file struct
		var honeyFileAttribs = fileAttribs{filepath, data.ModTime(), data.AccessTime()}
		return honeyFileAttribs
	}
	var honeyFileAttribs = fileAttribs{"", time.Now(), time.Now()}
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
func monitorHoneyFile(filepath string) {
	
	for {
		honeyAttribs1 := getFileAttribs(filepath)
		time.Sleep(2 * time.Second)
		honeyAttribs2 := getFileAttribs(filepath)
		fmt.Println(honeyAttribs1.accessTime)
		fmt.Println(honeyAttribs2.accessTime)
		if honeyAttribs1.modTime != honeyAttribs2.modTime || honeyAttribs1.accessTime != honeyAttribs2.accessTime {
			log.Fatal("[!] File has been accessed and/or modified! Sending alert!")
		} else {
			fmt.Println("[+] File untouched, resuming sleep...")
		}
	}
}

