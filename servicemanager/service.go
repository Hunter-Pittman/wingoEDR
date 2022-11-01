package servicemanager

import (
	"fmt"
	"unsafe"
	"golang.org/x/sys/windows"
)

// intializing the variables, calling the win32 apps dll ... helium is it's name
var (
	
	helium   = windows.NewLazyDLL("user32.dll")
    // how we can call functions
    accessSCM = helium.NewProc("GENERIC_READ")
	accessSrv = helium.NewProc("SC_MANAGER_ALL_ACCESS")


	hscm = helium.NewProc("OpenSCManager(None, None, accessSCM)")	



	serviceall = helium.NewProc("SERVICE_STATE_ALL")
	servicetype = helium.NewProc("SERVICE_WIN32")
	servicelist = helium.NewProc("EnumServicesStatusEx(hscm, servicetype, serviceall)")


)

func servicelister() {
	
	for (short_name, desc, status) in servicelist:
	print(short_name, desc, status) 
}

