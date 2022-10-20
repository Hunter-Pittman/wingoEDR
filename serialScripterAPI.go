package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"wingoEDR/common"

	"go.uber.org/zap"
)

type Beat struct {
	IP string
}

func HeartBeat() {
	m := Beat{IP: common.GetIP()}
	jsonStr, err := json.Marshal(m)
	if err != nil {
		zap.S().Warn(err)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Post("https://10.123.80.115:5000/api/v1/common/heartbeat", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	} else {
		println(resp.Status)
	}

	defer resp.Body.Close()

}

func IncidentAlert() {

}

func UpdateConfig() {

}
