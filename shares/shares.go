package shares

import (
	"fmt"
	"log"

	"github.com/Jeffail/gabs/v2"

	"github.com/abdfnx/gosh"
)

type SMBProperties struct {
	Description string
	Name        string
	Path        string
	Scoped      string
	ScopeName   string
	ShareState  string
	ShareType   string
	Special     string
	Temporary   string
	Namespace   string
	ServerName  string
	ClassName   string
}

func SmbShares() []SMBProperties {
	somatic := make([]SMBProperties, 0)

	err, out, errout := gosh.RunOutput("Get-SmbShare | ConvertTo-JSON")
	if err != nil {
		log.Printf("error: %v\n", err)
		fmt.Print(errout)
	}

	parent, err := gabs.ParseJSON([]byte(out))
	if err != nil {
		log.Printf("error: %v\n", err)
		fmt.Print(errout)
	}

	// firstgeneration := parent.Children()

	for _, child := range parent.Children() {
		chromosome := SMBProperties{
			Description: child.S("CimInstanceProperties").Children()[7].String(),
			Name:        child.S("CimInstanceProperties").Children()[15].String(),
			Path:        child.S("CimInstanceProperties").Children()[16].String(),
			Scoped:      child.S("CimInstanceProperties").Children()[19].String(),
			ScopeName:   child.S("CimInstanceProperties").Children()[20].String(),
			ShareState:  child.S("CimInstanceProperties").Children()[23].String(),
			Special:     child.S("CimInstanceProperties").Children()[26].String(),
			Temporary:   child.S("CimInstanceProperties").Children()[27].String(),
			Namespace:   child.Path("CimSystemProperties.Namespace").Data().(string),
			ServerName:  child.Path("CimSystemProperties.ServerName").Data().(string),
			ClassName:   child.Path("CimSystemProperties.ClassName").Data().(string)}

		somatic = append(somatic, chromosome)
		// for _, grandchild := range child.S("CimSystemProperties").Children() {
		// 	fmt.Println(grandchild.S("Namespace").String())

		// }

	}
	return (somatic)

}

// fmt.Println(sharess[1])

// var ss SharesSt
// err2 := json.Unmarshal([]byte(out), ss)
// if err2 != nil {
// 	log.Printf("error: %v\n", err2)
// 	fmt.Print(err2)
// }
// fmt.Println(out)
