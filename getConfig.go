package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Apis struct {
		KasperskyKey string `json:"kaspersky_key"`
	} `json:"apis"`
}

func getKaperskyKey() string {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration.Apis.KasperskyKey
}
