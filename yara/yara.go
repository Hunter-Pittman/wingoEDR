package yara

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strings"
	"wingoEDR/config"

	"go.uber.org/zap"
)

type YaraMatch struct {
	Rule string
	File string
}

// WARNING this will accpet directories of unlimited size. Directorires that are <= 1000 objects is recommended. Time estimate for a 1400 object directory would be 1 hour+
func YaraScan(rules string, dir string) ([]YaraMatch, error) {
	//var thing []YaraMatch
	ruleList := make([]string, 1)
	rulePath := rules

	yaraExePath := config.GetYaraExePath()

	files, err := ioutil.ReadDir(rulePath)
	if err != nil {
		zap.S().Warn(err)
	}

	for _, f := range files {
		ruleList = append(ruleList, f.Name())
	}

	matches := make([]YaraMatch, 0)
	ruleList = ruleList[1:]
	for i := range ruleList {
		cmd := exec.Command(yaraExePath, "-r", rulePath+ruleList[i], dir)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			zap.S().Fatal("cmd.Run() failed with %s\n", err)
		}
		outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		if errStr != "" {
			zap.S().Warn(errStr)
		}

		if outStr != "" {
			line := strings.SplitAfter(outStr, "\n")
			for x := range line {
				if line[x] != "" {
					rule := strings.SplitAfterN(line[x], " ", 2)
					match := YaraMatch{
						Rule: rule[0],
						File: rule[1],
					}
					matches = append(matches, match)
				}
			}
		}

	}
	return matches, nil
}
