package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"wingoEDR/common"

	"go.uber.org/zap"
)

var (
	client = http.Client{}
)

type HashLookup struct {
	Zone            string `json:"Zone"`
	FileGeneralInfo struct {
		FileStatus string    `json:"FileStatus"`
		Sha1       string    `json:"Sha1"`
		Md5        string    `json:"Md5"`
		Sha256     string    `json:"Sha256"`
		FirstSeen  time.Time `json:"FirstSeen"`
		LastSeen   time.Time `json:"LastSeen"`
		Signer     string    `json:"Signer"`
		Size       int       `json:"Size"`
		Type       string    `json:"Type"`
		HitsCount  int       `json:"HitsCount"`
	} `json:"FileGeneralInfo"`
}

// These are the possible output strings of this function (Malware, Adware and other, Clean, No threats detected, or Not categorized)
// Example calls
//getHashStatus("7a2278a9a74f49852a5d75c745ae56b80d5b4c16f3f6a7fdfd48cb4e2431c688") // Bad
//getHashStatus("408f31d86c6bf4a8aff4ea682ad002278f8cb39dc5f37b53d343e63a61f3cc4f") // Uncategorized
//getHashStatus("27dfb7631807c7bd185f57cd6de0628c6e9c47ed9b390a9b8544fdf12a323e04") // No results
//getHashStatus("98E07EDE313BAB4D2B659F4AF09804DB554287308EC1882D3D4036BEAE0D126E") // Clean

func GetHashStatus(hash string, hashType string) string {
	var isValidHash bool

	// Investigate to see if ToLower is needed for the hash switch case
	switch {
	case hashType == "md5":
		isValidHash = common.VerifyMD5Hash(hash)
	case hashType == "sha1":
		isValidHash = common.VerifySHA1Hash(hash)
	case hashType == "sha256":
		isValidHash = common.VerifySHA256Hash(hash)
	default:
		zap.S().Error("Supported hash types are: md5, sha1, sha256")
	}

	if !isValidHash {
		log.Fatalln("Hash verification failed!")
	}

	url := "https://opentip.kaspersky.com/api/v1/search/hash?request=" + hash

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("x-api-key", getKaperskyKey())
	response, err := client.Do(request)
	if err != nil {
		zap.S().Error("Kaspersky web request failed: ", err)
	}

	switch response.Status != "200 OK" {
	case "400 Bad Request" == response.Status:
		zap.S().Info("No results")
	case "401 Unauthorized" == response.Status:
		zap.S().Error("Your API key is not working!")
	case "403 Forbidden" == response.Status:
		zap.S().Error("You may have been blocked or your quote with kapersky reached: ", response.Status)
	case "404 Not Found" == response.Status:
		zap.S().Error("Resource not found, check your endpoint address: ", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		zap.S().Error(err)
	}

	var result HashLookup

	sb := string(body)
	zap.S().Info(sb)

	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		zap.S().Error("Cannot unmarshal JSON")
	}

	return result.FileGeneralInfo.FileStatus
}
