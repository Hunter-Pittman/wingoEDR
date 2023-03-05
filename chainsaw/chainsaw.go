package chainsaw

import (
	"os/exec"
	"regexp"
	"strings"
	"wingoEDR/common"
	"wingoEDR/config"

	"github.com/Jeffail/gabs"
	"github.com/fatih/color"
	"go.uber.org/zap"
)

var (
	CHAINSAW_PATH         = config.GetChainsawPath()
	CHAINSAW_MAPPING_PATH = config.GetChainsawMapping()
	CHAINSAW_RULES_PATH   = config.GetChainSawRulesBad()
)

func ScanAll() *gabs.Container {
	//cmdOutput, err := exec.Command("c:\\Users\\Hunter Pittman\\Documents\\repos\\wingoEDR\\externalresources\\chainsaw\\chainsaw.exe", "--no-banner", "hunt", "C:\\Windows\\System32\\winevt\\Logs", "-s", "C:\\Users\\Hunter Pittman\\Documents\\repos\\wingoEDR\\externalresources\\chainsaw\\rules\\bad", "--mapping", "C:\\Users\\Hunter Pittman\\Documents\\repos\\wingoEDR\\externalresources\\chainsaw\\mappings\\sigma-event-logs-all.yml", "--json").Output()
	cmdOutput, err := exec.Command(CHAINSAW_PATH, "--no-banner", "hunt", "C:\\Windows\\System32\\winevt\\Logs\\", "-s", CHAINSAW_RULES_PATH, "--mapping", CHAINSAW_MAPPING_PATH, "--json").Output()
	if err != nil {
		if strings.Contains(err.Error(), "Specified event log path is invalid") {
			zap.S().Error("Error opening evtx log files: ", err.Error())
			color.Red("[ERROR]	Failed opening evtx log files: ", err.Error())
			return nil
		}
	} else {
		common.ErrorHandler(err)
	}

	common.ErrorHandler(err)

	parsedJSON, err := gabs.ParseJSON(cmdOutput)

	//println(string(cmdOutput))

	return parsedJSON
}

func ScanTimeRange(fromTimestamp string, toTimestamp string) *gabs.Container {
	// Example: --from "2023-03-04T00:00:00" --to "2023-03-05T23:59:59"

	timestampPattern := "[0-9]{4}-[0-9]{2}-[0-9]{2}T(0?[0-9]|1[0-9]|2[0-3]):[0-9]+:(0?[0-9]|[1-5][0-9])"

	match, _ := regexp.MatchString(timestampPattern, fromTimestamp)

	if match == false {
		color.Red("[ERROR]	Invalid timestamp format")
		return nil
	}

	match, _ = regexp.MatchString(timestampPattern, toTimestamp)

	if match == false {
		color.Red("[ERROR]	Invalid timestamp format")
		return nil
	}

	fromTimestamp = common.LocalTimeToUTC(fromTimestamp)
	toTimestamp = common.LocalTimeToUTC(toTimestamp)

	cmdOutput, err := exec.Command(CHAINSAW_PATH, "--no-banner", "hunt", "C:\\Windows\\System32\\winevt\\Logs\\", "-s", CHAINSAW_RULES_PATH, "--mapping", CHAINSAW_MAPPING_PATH, "--from", fromTimestamp, "--to", toTimestamp, "--json").Output()
	if err != nil {
		if strings.Contains(err.Error(), "Specified event log path is invalid") {
			zap.S().Error("Error opening evtx log files: ", err.Error())
			color.Red("[ERROR]	Failed opening evtx log files: ", err.Error())
			return nil
		}
	} else {
		common.ErrorHandler(err)
	}

	parsedJSON, err := gabs.ParseJSON(cmdOutput)

	//println(string(cmdOutput))

	return parsedJSON
}
