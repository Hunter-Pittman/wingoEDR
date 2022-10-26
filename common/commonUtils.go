package common

import (
	"net"
	"regexp"
	"time"
	"github.com/djherbis/times"
	"go.uber.org/zap"
	"strings"
	"crypto/sha1"
	"encoding/hex"
)

type fileAttribs struct {
	filename	string
	modTime		time.Time
	accessTime	time.Time
}

func GetIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		zap.S().Warn(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	ipaddr := localAddr.IP
	return ipaddr.String()
}

func GenerateSha1Hash(data string) string {
	dataHashByte := sha1.Sum([]byte(data))
	dataHashStr := hex.EncodeToString(dataHashByte[:])
	return dataHashStr
}

func VerifySHA256Hash(hash string) bool {
	match, _ := regexp.MatchString("[A-Fa-f0-9]{64}", hash)
	return match
}

func getFileAttribs(filepath string) fileAttribs {
	data, err := times.Stat(filepath)
	if err != nil {
		if strings.Contains(err.Error(), "cannot find the file") {
			// sign that the file was deleted
			zap.S().Warn("1 honey file was likely deleted! Sending alert:", err)
		}
	} else {
		// create file struct
		var honeyDirAttribs = fileAttribs{filepath, data.ModTime(), data.AccessTime()}
		return honeyDirAttribs
	}
	var honeyFileAttribs = fileAttribs{"", time.Now(), time.Now()}
	return honeyFileAttribs

}


func VerifySHA1Hash(hash string) bool {
	match, _ := regexp.MatchString("[a-fA-F0-9]{40}$", hash)
	return match
}

func VerifyMD5Hash(hash string) bool {
	match, _ := regexp.MatchString("/^[a-f0-9]{32}$/i", hash)
	return match
}

func FirstWords(value string, count int) string {
	// Loop over all indexes in the string.
	for i := range value {
		// If we encounter a space, reduce the count.
		if value[i] == ' ' {
			count -= 1
			// When no more words required, return a substring.
			if count == 0 {
				return value[0:i]
			}
		}
	}
	// Return the entire string.
	return value
}
