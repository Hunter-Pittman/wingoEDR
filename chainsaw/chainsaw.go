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

type ChainsawEvent struct {
	ID         string
	RuleName   string
	Timestamp  string
	Tags       []string
	Authors    []string
	Level      string
	References []string
	Payload    map[string]interface{}
}

func ScanAll() ([]ChainsawEvent, error) {
	var CHAINSAW_PATH string = config.GetChainsawPath()
	var CHAINSAW_MAPPING_PATH string = config.GetChainsawMapping()
	var CHAINSAW_RULES_PATH string = config.GetChainSawRulesBad()
	cmdOutput, err := exec.Command(CHAINSAW_PATH, "--no-banner", "hunt", "C:\\Windows\\System32\\winevt\\Logs\\", "-s", CHAINSAW_RULES_PATH, "--mapping", CHAINSAW_MAPPING_PATH, "--json").Output()
	if err != nil {
		if strings.Contains(err.Error(), "Specified Chainsawevent log path is invalid") {
			zap.S().Error("Error opening evtx log files: ", err.Error())
			color.Red("[ERROR]	Failed opening evtx log files: ", err.Error())
			return nil, err
		} else {
			zap.S().Error("Error running chainsaw: ", err.Error())
		}
	}

	parsedJSON, err := gabs.ParseJSON(cmdOutput)
	if err != nil {
		zap.S().Error("Error parsing JSON: ", err.Error())
		return nil, err
	}
	newChainsawEvent, err := chainsawToStruct(parsedJSON)
	if err != nil {
		zap.S().Error("Error converting chainsaw output to struct: ", err.Error())
		return nil, err
	}

	//println(string(cmdOutput))

	return newChainsawEvent, nil
}

func ScanAllJSON() (*gabs.Container, error) {
	var CHAINSAW_PATH string = config.GetChainsawPath()
	var CHAINSAW_MAPPING_PATH string = config.GetChainsawMapping()
	var CHAINSAW_RULES_PATH string = config.GetChainSawRulesBad()
	cmdOutput, err := exec.Command(CHAINSAW_PATH, "--no-banner", "hunt", "C:\\Windows\\System32\\winevt\\Logs\\", "-s", CHAINSAW_RULES_PATH, "--mapping", CHAINSAW_MAPPING_PATH, "--json").Output()
	if err != nil {
		if strings.Contains(err.Error(), "Specified Chainsawevent log path is invalid") {
			zap.S().Error("Error opening evtx log files: ", err.Error())
			color.Red("[ERROR]	Failed opening evtx log files: ", err.Error())
			return nil, err
		} else {
			zap.S().Error("Error encountered with chainsaw.exe: ", err.Error())
		}
	}

	parsedJSON, err := gabs.ParseJSON(cmdOutput)
	if err != nil {
		zap.S().Error("Error parsing JSON: ", err.Error())
		return nil, err
	}

	//println(string(cmdOutput))

	return parsedJSON, nil
}

func ScanTimeRange(fromTimestamp string, toTimestamp string) ([]ChainsawEvent, error) {
	var CHAINSAW_PATH string = config.GetChainsawPath()
	var CHAINSAW_MAPPING_PATH string = config.GetChainsawMapping()
	var CHAINSAW_RULES_PATH string = config.GetChainSawRulesBad()
	// Example: --from "2023-03-04T00:00:00" --to "2023-03-05T23:59:59"

	timestampPattern := `(0?[1-9]|[1][0-2])-[0-9]+-[0-9]+T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]{1,3})?`

	match, _ := regexp.MatchString(timestampPattern, fromTimestamp)

	if !match {
		zap.S().Error("Invalid timestamp format for: --from")
		return nil, errors.New("Invalid timestamp format for from")
	}

	match, _ = regexp.MatchString(timestampPattern, toTimestamp)

	if !match {
		zap.S().Error("Invalid timestamp format for: --to")
		return nil, errors.New("Invalid timestamp format for to")
	}

	fromTimestamp, err := common.LocalTimeToUTC(fromTimestamp)
	if err != nil {
		zap.S().Error("Error converting fromTimestamp to UTC: ", err.Error())
		return nil, err
	}

	toTimestamp, err1 := common.LocalTimeToUTC(toTimestamp)
	if err1 != nil {
		zap.S().Error("Error converting toTimestamp to UTC: ", err1.Error())
		return nil, err1
	}

	cmdOutput, err := exec.Command(CHAINSAW_PATH, "--no-banner", "hunt", "C:\\Windows\\System32\\winevt\\Logs\\", "-s", CHAINSAW_RULES_PATH, "--mapping", CHAINSAW_MAPPING_PATH, "--from", fromTimestamp, "--to", toTimestamp, "--json").Output()
	if err != nil {
		if strings.Contains(err.Error(), "Specified Chainsawevent log path is invalid") {
			zap.S().Error("Error opening evtx log files: ", err.Error())
			color.Red("[ERROR]	Failed opening evtx log files: ", err.Error())
			return nil, err
		} else {
			zap.S().Error("Error encountered with chainsaw.exe: ", err.Error())
		}
	}

	parsedJSON, err := gabs.ParseJSON(cmdOutput)
	if err != nil {
		zap.S().Error("Error parsing JSON: ", err.Error())
		return nil, err
	}
	newChainsawEvent, err := chainsawToStruct(parsedJSON)
	if err != nil {
		zap.S().Error("Error converting chainsaw output to struct: ", err.Error())
		return nil, err
	}

	//println(string(cmdOutput))

	return newChainsawEvent, nil
}

func ScanTimeRangeJSON(fromTimestamp string, toTimestamp string) (*gabs.Container, error) {
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

	fromTimestamp, err := common.LocalTimeToUTC(fromTimestamp)
	if err != nil {
		zap.S().Error("Error converting fromTimestamp to UTC: ", err.Error())
		return nil, err
	}
	toTimestamp, err1 := common.LocalTimeToUTC(toTimestamp)
	if err1 != nil {
		zap.S().Error("Error converting toTimestamp to UTC: ", err1.Error())
		return nil, err1
	}

	cmdOutput, err := exec.Command(CHAINSAW_PATH, "--no-banner", "hunt", "C:\\Windows\\System32\\winevt\\Logs\\", "-s", CHAINSAW_RULES_PATH, "--mapping", CHAINSAW_MAPPING_PATH, "--from", fromTimestamp, "--to", toTimestamp, "--json").Output()
	if err != nil {
		if strings.Contains(err.Error(), "Specified Chainsawevent log path is invalid") {
			zap.S().Error("Error opening evtx log files: ", err.Error())
			color.Red("[ERROR]	Failed opening evtx log files: ", err.Error())
			return nil, err
		} else {
			zap.S().Error("Error encountered with chainsaw.exe: ", err.Error())
		}
	}

	parsedJSON, err := gabs.ParseJSON(cmdOutput)
	if err != nil {
		zap.S().Error("Error parsing JSON: ", err.Error())
		return nil, err
	}

	//println(string(cmdOutput))

	return parsedJSON, nil
}

func chainsawToStruct(chainsawOutput *gabs.Container) ([]ChainsawEvent, error) {
	var Chainsawevents []ChainsawEvent

	children, _ := chainsawOutput.Children()

	for _, child := range children {
		var Chainsawevent ChainsawEvent
		Chainsawevent.ID = child.Path("id").Data().(string)
		Chainsawevent.RuleName = child.Path("name").Data().(string)
		Chainsawevent.Level = child.Path("level").Data().(string)

		// Convert UTC timestamp to local timestamp and add to struct
		timestamp := child.Path("timestamp").Data().(string)
		utcToLocalTimestamp, err := common.UTCToLocalTime(timestamp)
		if err != nil {
			zap.S().Error("Error converting UTC timestamp to local timestamp: ", err.Error())
			return nil, err
		}
		Chainsawevent.Timestamp = utcToLocalTimestamp

		// Payload may or may not have a value
		//Chainsawevent.Payload = child.Path("document.data").Data().(map[string]interface{})

		if payload := child.Path("document.data").Data(); payload != nil {
			if payloadMap, ok := payload.(map[string]interface{}); ok {
				Chainsawevent.Payload = payloadMap
			} else {
				Chainsawevent.Payload = nil
			}
		}

		// Tags
		tags, _ := child.Path("tags").Children()
		for _, tag := range tags {
			Chainsawevent.Tags = append(Chainsawevent.Tags, tag.Data().(string))
		}

		// Authors
		authors, _ := child.Path("authors").Children()
		for _, author := range authors {
			Chainsawevent.Authors = append(Chainsawevent.Authors, author.Data().(string))
		}

		// References
		references, _ := child.Path("references").Children()
		for _, reference := range references {
			Chainsawevent.References = append(Chainsawevent.References, reference.Data().(string))
		}

		Chainsawevents = append(Chainsawevents, Chainsawevent)
	}

	return Chainsawevents, nil
}
