package chainsaw

import (
	"errors"
	"os/exec"
	"regexp"
	"strings"
	"wingoEDR/common"
	"wingoEDR/config"

	"github.com/Jeffail/gabs"
	"github.com/fatih/color"
	"go.uber.org/zap"
)

type Event struct {
	EventID  string
	RuleName string
	Payload  []byte
}

func ScanAll() (*gabs.Container, error) {
	var CHAINSAW_PATH string = config.GetChainsawPath()
	var CHAINSAW_MAPPING_PATH string = config.GetChainsawMapping()
	var CHAINSAW_RULES_PATH string = config.GetChainSawRulesBad()
	//cmdOutput, err := exec.Command("c:\\Users\\Hunter Pittman\\Documents\\repos\\wingoEDR\\externalresources\\chainsaw\\chainsaw.exe", "--no-banner", "hunt", "C:\\Windows\\System32\\winevt\\Logs", "-s", "C:\\Users\\Hunter Pittman\\Documents\\repos\\wingoEDR\\externalresources\\chainsaw\\rules\\bad", "--mapping", "C:\\Users\\Hunter Pittman\\Documents\\repos\\wingoEDR\\externalresources\\chainsaw\\mappings\\sigma-event-logs-all.yml", "--json").Output()
	cmdOutput, err := exec.Command(CHAINSAW_PATH, "--no-banner", "hunt", "C:\\Windows\\System32\\winevt\\Logs\\", "-s", CHAINSAW_RULES_PATH, "--mapping", CHAINSAW_MAPPING_PATH, "--json").Output()
	if err != nil {
		if strings.Contains(err.Error(), "Specified event log path is invalid") {
			zap.S().Error("Error opening evtx log files: ", err.Error())
			color.Red("[ERROR]	Failed opening evtx log files: ", err.Error())
			return nil, err
		} else {
			common.ErrorHandler(err)
		}
	}

	parsedJSON, err := gabs.ParseJSON(cmdOutput)

	//println(string(cmdOutput))

	return parsedJSON, nil
}

func ScanTimeRange(fromTimestamp string, toTimestamp string) (*gabs.Container, error) {
	var CHAINSAW_PATH string = config.GetChainsawPath()
	var CHAINSAW_MAPPING_PATH string = config.GetChainsawMapping()
	var CHAINSAW_RULES_PATH string = config.GetChainSawRulesBad()
	// Example: --from "2023-03-04T00:00:00" --to "2023-03-05T23:59:59"

	timestampPattern := `(0?[1-9]|[1][0-2])-[0-9]+-[0-9]+T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]{1,3})?`

	match, _ := regexp.MatchString(timestampPattern, fromTimestamp)

	if !match {
		color.Red("[ERROR]	Invalid timestamp format for: --from")
		return nil, errors.New("Invalid timestamp format for from")
	}

	match, _ = regexp.MatchString(timestampPattern, toTimestamp)

	if !match {
		color.Red("[ERROR]	Invalid timestamp format for: --to")
		return nil, errors.New("Invalid timestamp format for from")
	}

	fromTimestamp = common.LocalTimeToUTC(fromTimestamp)
	toTimestamp = common.LocalTimeToUTC(toTimestamp)

	cmdOutput, err := exec.Command(CHAINSAW_PATH, "--no-banner", "hunt", "C:\\Windows\\System32\\winevt\\Logs\\", "-s", CHAINSAW_RULES_PATH, "--mapping", CHAINSAW_MAPPING_PATH, "--from", fromTimestamp, "--to", toTimestamp, "--json").Output()
	if err != nil {
		if strings.Contains(err.Error(), "Specified event log path is invalid") {
			zap.S().Error("Error opening evtx log files: ", err.Error())
			color.Red("[ERROR]	Failed opening evtx log files: ", err.Error())
			return nil, err
		} else {
			common.ErrorHandler(err)
		}
	}

	parsedJSON, err := gabs.ParseJSON(cmdOutput)

	//println(string(cmdOutput))

	return parsedJSON, nil
}
