package shares

import (
	"fmt"

	win "github.com/gorpher/gowin32"
	"go.uber.org/zap"
)

type SMBINFO struct {
	Netname     string
	Remark      string
	Path        string
	Type        int
	Permissions int
	MaxUses     int
	CurrentUses int
}

func shares() []SMBINFO {
	x, err := fmt.Println("")
	if err != nil {
		//log.Fatal(err)
		zap.Error(err)
	}
	fmt.Println(x)
	shareslice := make([]SMBINFO, 0)

	share := win.NetShareEnum()

	for _, v := range share {

		helium := SMBINFO{
			Netname:     v.Netname,
			Remark:      v.Remark,
			Path:        v.Path,
			Type:        v.Permissions,
			MaxUses:     v.MaxUses,
			CurrentUses: v.CurrentUses}

		shareslice = append(shareslice, helium)

	}
	return (shareslice)

}
