package shares

import (
	win "github.com/gorpher/gowin32"
)

type SMBInfo struct {
	Netname     string
	Remark      string
	Path        string
	Type        int
	Permissions int
	MaxUses     int
	CurrentUses int
}

func GetShares() []SMBInfo {

	shareslice := make([]SMBInfo, 0)

	share := win.NetShareEnum()

	for _, v := range share {

		helium := SMBInfo{
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
