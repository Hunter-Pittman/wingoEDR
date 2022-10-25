package yara

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

type YaraMatch struct {
	rule string
	file string
}

func DirYaraScan() {
	ruleList := make([]string, 1)
	rulePath := "C:\\Users\\hunte\\Documents\\repos\\wingoEDR\\yara_rules\\fileID\\"

	files, err := ioutil.ReadDir(rulePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		ruleList = append(ruleList, f.Name())
	}

	//var yaraMatches []YaraMatch
	ruleList = ruleList[1:]
	for i := range ruleList {
		cmd := exec.Command("C:\\Users\\hunte\\Documents\\repos\\wingoEDR\\yara_exe\\yara64.exe", "-r", rulePath+ruleList[i], "C:\\Users\\hunte\\Pictures")

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			zap.S().Fatal("cmd.Run() failed with %s\n", err)
		}
		outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		if errStr != "" {
			zap.S().Warn("There is an error", errStr)
		}

		if outStr != "" {
			line := strings.SplitAfter(outStr, "\n")
			for x := range line {
				rule := strings.SplitAfterN(line[x], " ", 2)
				fmt.Printf("%v\n", rule[0])
			}
		}

	}
}

func FileYaraScan() {

}
