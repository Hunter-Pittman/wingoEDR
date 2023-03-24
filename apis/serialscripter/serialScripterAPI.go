package serialscripter

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"wingoEDR/config"
	"wingoEDR/inventory"

	"go.uber.org/zap"
)

//https://ec2-18-246-47-205.us-west-2.compute.amazonaws.com:10000

type Beat struct {
	IP string
}
type Incident struct {
	Name        string
	CurrentTime string
	User        string
	Severity    string
	Payload     string
}

type Alert struct {
	Host     string
	Incident Incident
}

func HeartBeat() (err error) {
	runRequest := CheckEndpoint(false)
	if !runRequest {
		return nil
	}

	ssUserAgent := config.GetSerialScripterUserAgent()
	//ssUserAgent := "nestler-code"

	if err != nil {
		zap.S().Warn(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	requestURL := fmt.Sprintf("%v/api/v1/common/heartbeat", config.GetSerialScripterURL())
	req, err := http.NewRequest(http.MethodPost, requestURL, nil)
	if err != nil {
		zap.S().Warn(err)
	}

	req.Header.Set("User-Agent", ssUserAgent)
	resp, err := client.Do(req)
	if err != nil {
		zap.S().Error(err)
		return err
	} else {
		// data, _ := ioutil.ReadAll(resp.Body)
		// println(string(data))
	}

	defer resp.Body.Close()
	return nil

}

func PostInventory() (err error) {
	runRequest := CheckEndpoint(false)
	if !runRequest {
		return nil
	}
	ssUserAgent := config.GetSerialScripterUserAgent()

	// Payload
	inventoryItems := inventory.GetInventory()

	jsonStr, err := json.Marshal(inventoryItems)
	if err != nil {
		zap.S().Warn(err)
	}

	//fmt.Println(string(jsonStr))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	bodyReader := bytes.NewReader(jsonStr)

	requestURL := fmt.Sprintf("%v/api/v1/common/inventory", config.GetSerialScripterURL())
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		zap.S().Warn(err)
	}

	req.Header.Set("User-Agent", ssUserAgent)
	resp, err := client.Do(req)
	if err != nil {
		zap.S().Error(err)
		return err
	} else {
		// data, _ := ioutil.ReadAll(resp.Body)
		// println(string(data))
	}

	defer resp.Body.Close()
	return nil
}

func IncidentAlert(alert Alert) (err error) {
	runRequest := CheckEndpoint(false)
	if !runRequest {
		return nil
	}

	ssUserAgent := config.GetSerialScripterUserAgent()

	jsonStr, err := json.Marshal(alert)
	if err != nil {
		zap.S().Warn(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	bodyReader := bytes.NewReader(jsonStr)

	requestURL := fmt.Sprintf("%v/api/v1/common/incidentalert", config.GetSerialScripterURL())
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		zap.S().Warn(err)
	}

	req.Header.Set("User-Agent", ssUserAgent)
	resp, err := client.Do(req)
	if err != nil {
		zap.S().Error(err)
		return err
	} else {
		// data, _ := ioutil.ReadAll(resp.Body)
		// println(string(data))
	}

	defer resp.Body.Close()
	return nil
}

func CheckEndpoint(flag bool) bool {
	url := config.GetSerialScripterURL()
	if flag {
		return false
	} else if !strings.Contains(url, "http") {
		return false
	} else if TestHeartBeat() != nil {
		return false
	} else {
		return true
	}
}

func TestHeartBeat() (err error) {

	ssUserAgent := config.GetSerialScripterUserAgent()
	//ssUserAgent := "nestler-code"

	if err != nil {
		zap.S().Warn(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	requestURL := fmt.Sprintf("%v/api/v1/common/heartbeat", config.GetSerialScripterURL())
	req, err := http.NewRequest(http.MethodPost, requestURL, nil)
	if err != nil {
		zap.S().Warn(err)
	}

	req.Header.Set("User-Agent", ssUserAgent)
	resp, err := client.Do(req)
	if err != nil {
		zap.S().Error(err)
		return err
	} else {
		if resp.StatusCode != http.StatusOK {
			zap.S().Errorf("Unexpected status code %d\n", resp.StatusCode)
			return errors.New("Unexpected status code")
		}
	}

	defer resp.Body.Close()
	return nil

}
