package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Beat struct {
	IP string
}

func HeartBeat() {
	m := Beat{IP: GetIP()}
	jsonStr, err := json.Marshal(m)
	if err != nil {
		zap.S().Warn(err)
	}
	resp, err := http.Post("http://localhost:80/heartbeat", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		panic(err)
	} else {
		println(resp.Status)
	}

	defer resp.Body.Close()
}
