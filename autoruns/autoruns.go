package autoruns

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
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

const AUTORUNS_LOC = "C:\\Users\\hunte\\Documents\\repos\\wingoEDR\\externalbinaries\\autorunsc.exe"

func FullAutorunsDump() {

	xmlName := strconv.FormatInt(time.Now().Unix(), 10) + "_autoruns.xml"
	fmt.Println(xmlName)

	cmdOutput, err := exec.Command(AUTORUNS_LOC, "-a", "*", "-x").Output()
	if err != nil {
		log.Fatal(err)
	}

	file, e := os.Create(xmlName)
	if e != nil {
		log.Fatal(err)
	}

	stringifiedCmdOutput := string(cmdOutput)[275:]

	_, err2 := file.WriteString(string(stringifiedCmdOutput))
	if err2 != nil {
		log.Fatal(err)
	}

	byteValue, _ := ioutil.ReadAll(file)

	e1 := &Start{}
	err = xml.Unmarshal([]byte(byteValue), &e1)
	if err != nil {
		log.Fatal(err)
	}
}
