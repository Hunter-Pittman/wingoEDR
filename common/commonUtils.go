package common

import (
	"bufio"
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/djherbis/times"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/unicode"
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

func VerifyWindowsPathFatal(path string) {
	match, _ := regexp.MatchString(`[a-zA-Z]:[\\\/](?:[a-zA-Z0-9]+[\\\/])*([a-zA-Z0-9]+\.*)`, path)

	if !match { // Errors out on a "C:\" path needs to be fixed
		color.Red("[ERROR]	The entered output is not a Windows path!")
		os.Exit(1)
	} else {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			color.Red("[ERROR]	Windows path does not exist!")
			os.Exit(1)
		}
	}
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

func GetSerialScripterHostName() string {
	lastOctets := strings.Split(GetIP(), ".")
	serialScripterHostName := "host-" + lastOctets[3]

	return serialScripterHostName
}

func CsvToJsonSysInternals(csvFile string) (string, error) {
	// Open the CSV file
	f, err := os.Open(csvFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Decode UTF-16 encoded CSV file
	utf16Decoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	reader := utf16Decoder.Reader(f)

	// Read the CSV file
	r := csv.NewReader(bufio.NewReader(reader))
	r.Comma = ';'
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		return "", err
	}

	// Convert CSV records to a slice of maps
	var result []map[string]string
	header := records[0]
	for _, row := range records[1:] {
		record := make(map[string]string)
		for i, v := range row {
			if i < len(header) {
				record[header[i]] = v
			} else {
				record[fmt.Sprintf("field%d", i+1)] = v
			}
		}
		result = append(result, record)
	}

	// Convert the slice of maps to JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	// Return the JSON as a string
	return string(jsonData), nil
}

func ErrorHandler(err error) {
	if err != nil {
		zap.S().Error("Error: ", err.Error())
		color.Red("[ERROR]	An error has been encounterd: ", err.Error())
	} else {
		return
	}
}
