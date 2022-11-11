package common

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/djherbis/times"
	"go.uber.org/zap"
)

type FileAttribs struct {
	filename   string
	modTime    time.Time
	accessTime time.Time
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

func GetFileAttribs(filepath string) FileAttribs {
	data, err := times.Stat(filepath)
	if err != nil {
		if strings.Contains(err.Error(), "cannot find the file") {
			// sign that the file was deleted
			zap.S().Warn("1 honey file was likely deleted! Sending alert:", err)
		}
	} else {
		// create file struct
		var honeyDirAttribs = FileAttribs{filepath, data.ModTime(), data.AccessTime()}
		return honeyDirAttribs
	}
	var honeyFileAttribs = FileAttribs{"", time.Now(), time.Now()}
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

type finfo struct {
	Name string
	Size int64
	Time string
	Hash string
}

func CheckFile(name string) finfo {
	fileInfo, err := os.Stat(name)
	if err != nil {
		panic(err)
	}
	println(name)
	if fileInfo.IsDir() {

		t := fileInfo.ModTime().String()
		b := fileInfo.Size()

		i := finfo{
			Name: name,
			Size: b,
			Time: t,
			Hash: "directory",
		}

		return i
	} else {
		f, err := os.Open(name)
		if err != nil {
			zap.S().Error(err)
		}
		if err != nil {
			if os.IsNotExist(err) {
				println("file not found:", fileInfo.Name())
			}
		}
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			zap.S().Warn(err)
		}
		hash := h.Sum(nil)
		Enc := base64.StdEncoding.EncodeToString(hash)

		t := fileInfo.ModTime().String()
		b := fileInfo.Size()

		i := finfo{
			Name: name,
			Size: b,
			Time: t,
			Hash: Enc,
		}
		return i
	}
}

func GetSerialScripterHostName() string {
	lastOctets := strings.Split(GetIP(), ".")
	serialScripterHostName := "host-" + lastOctets[3]

	return serialScripterHostName
}

// Repalce with a DC query????
func GetCurrentlyLoggedInUsers() string {
	toExecute := "query user /server:$SERVER | ConvertTo-Json"
	output, err := exec.Command("powershell.exe", toExecute).Output()
	if err != nil {
		zap.S().Error(err.Error())
		//fmt.Println(err)
	}

	return string(output)
}
