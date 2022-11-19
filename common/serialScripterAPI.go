package common

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Beat struct {
	IP string
}
type Incident struct {
	Name     string
	User     string
	Process  string
	RemoteIP string
	Cmd      string
}

type Alert struct {
	Host     string
	Incident Incident
}

func HeartBeat() {
	ssUserAgent := GetSerialScripterUserAgent()

	m := Beat{IP: GetIP()}
	jsonStr, err := json.Marshal(m)
	if err != nil {
		zap.S().Warn(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	bodyReader := bytes.NewReader(jsonStr)

	requestURL := fmt.Sprintf("https://10.123.80.115:10000/api/v1/common/heartbeat")
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		zap.S().Warn(err)
	}

	req.Header.Set("User-Agent", ssUserAgent)
	resp, err := client.Do(req)
	if err != nil {
		zap.S().Error(err)
	} else {
		//data, _ := ioutil.ReadAll(resp.Body)
		//println(string(data))
	}

	defer resp.Body.Close()

}

func PostInventory() {
	ssUserAgent := GetSerialScripterUserAgent()
	inventoryItems := GetInventory()

	jsonStr, err := json.Marshal(inventoryItems)
	if err != nil {
		zap.S().Warn(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	bodyReader := bytes.NewReader(jsonStr)

	requestURL := fmt.Sprintf("https://10.123.80.115:10000/api/v1/common/inventory")
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		zap.S().Warn(err)
	}

	req.Header.Set("User-Agent", ssUserAgent)
	resp, err := client.Do(req)
	if err != nil {
		zap.S().Error(err)
	} else {
		//data, _ := ioutil.ReadAll(resp.Body)
		//println(string(data))
	}

	defer resp.Body.Close()
}

func IncidentAlert(alert Alert) {
	ssUserAgent := GetSerialScripterUserAgent()

	jsonStr, err := json.Marshal(alert)
	if err != nil {
		zap.S().Warn(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	bodyReader := bytes.NewReader(jsonStr)

	requestURL := fmt.Sprintf("https://10.123.80.115:10000/api/v1/common/incidentalert")
	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		zap.S().Warn(err)
	}

	req.Header.Set("User-Agent", ssUserAgent)
	resp, err := client.Do(req)
	if err != nil {
		zap.S().Error(err)
	} else {
		//data, _ := ioutil.ReadAll(resp.Body)
		//println(string(data))
	}

	defer resp.Body.Close()
}