package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
)

type Beat struct {
	IP string
}

func HeartBeat() {
	m := Beat{IP: getIP()}
	jsonStr, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	resp, err := http.Post("http://localhost:80/heartbeat", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		panic(err)
	} else {
		println(resp.Status)
	}

	defer resp.Body.Close()
}

func getIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	ipaddr := localAddr.IP
	return ipaddr.String()
}
