package autoruns

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
	"wingoEDR/common"
)

type Start struct {
	Autoruns Autoruns `json:"autoruns"`
}

type Autoruns struct {
	Item AutorunsEntry `json:"item"`
}

type AutorunsEntry struct {
	Location     string `json:"location"`
	Itemname     string `json:"itemname"`
	Enabled      string `json:"enabled"`
	Profile      string `json:"profile"`
	Launchstring string `json:"launchstring"`
	Description  string `json:"description"`
	Company      string `json:"company"`
	Version      string `json:"version"`
	Imagepath    string `json:"imagepath"`
	Time         string `json:"time"`
}

const AUTORUNS_LOC = "C:\\Users\\FORENSICS\\Documents\\Hunter's Repos\\wingoEDR\\externalbinaries\\autorunsc.exe"

func FullAutorunsDump() {

	csvName := strconv.FormatInt(time.Now().Unix(), 10) + "_autoruns.csv"
	fmt.Println(csvName)

	cmdOutput, err := exec.Command(AUTORUNS_LOC, "-nobanner", "-a", "*", "-c").Output()
	if err != nil {
		log.Fatal(err)
	}

	file, e := os.Create(csvName)
	if e != nil {
		log.Fatal(err)
	}

	stringifiedCmdOutput := string(cmdOutput)

	_, err2 := file.WriteString(string(stringifiedCmdOutput))
	if err2 != nil {
		log.Fatal(err)
	}

	json, err := common.CsvToJsonSysInternals(csvName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(json)
}
