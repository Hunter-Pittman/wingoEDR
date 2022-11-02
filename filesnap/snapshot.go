package filesnap


import (
	"wingoEDR/honeytoken"
)


func GatherAttributes(filepath string) honeytoken.FileAttribs {
	attribs := honeytoken.GetFileAttribs(filepath)
	return attribs
}