package filesnap

import (
	"wingoEDR/honeymonitor"
)

func GatherAttributes(filepath string) honeymonitor.FileAttribs {
	attribs := honeymonitor.GetFileAttribs(filepath)
	return attribs
}
