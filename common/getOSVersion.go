package common

import (
	"fmt"
	"syscall"
	"unsafe"
)

type osVersionInfoEx struct {
	dwOSVersionInfoSize uint32
	dwMajorVersion      uint32
	dwMinorVersion      uint32
	dwBuildNumber       uint32
	dwPlatformId        uint32
	szCSDVersion        [128]uint16
	wServicePackMajor   uint16
	wServicePackMinor   uint16
	wSuiteMask          uint16
	wProductType        byte
	wReserved           byte
}

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var getVersionEx = kernel32.NewProc("GetVersionExW")

func getWindowsVersion() (string, error) {
	var info osVersionInfoEx
	info.dwOSVersionInfoSize = uint32(unsafe.Sizeof(info))

	ret, _, err := getVersionEx.Call(uintptr(unsafe.Pointer(&info)))
	if ret == 0 {
		return "", err
	}

	version := fmt.Sprintf("Windows %d.%d", info.dwMajorVersion, info.dwMinorVersion)
	if info.dwMajorVersion == 10 {
		switch info.dwMinorVersion {
		case 0:
			version += " (Windows 10)"
		case 1:
			version += " (Windows 10 Anniversary Update)"
		case 2:
			version += " (Windows 10 Creators Update)"
		case 3:
			version += " (Windows 10 Fall Creators Update)"
		case 4:
			version += " (Windows 10 April 2018 Update)"
		case 5:
			version += " (Windows 10 October 2018 Update)"
		case 6:
			version += " (Windows 10 May 2019 Update)"
		case 7:
			version += " (Windows 10 November 2019 Update)"
		case 8:
			version += " (Windows 10 May 2020 Update)"
		case 9:
			version += " (Windows 10 October 2020 Update)"
		case 10:
			version += " (Windows 10 May 2021 Update)"
		default:
			version += fmt.Sprintf(" (%d.%d)", info.dwMajorVersion, info.dwMinorVersion)
		}
	}

	if info.dwMajorVersion == 6 {
		switch info.dwMinorVersion {
		case 0:
			version += " (Windows Vista)"
		case 1:
			version += " (Windows 7)"
		case 2:
			if info.wServicePackMajor == 0 {
				version += " (Windows 8)"
			} else {
				version += " (Windows 8.1)"
			}
		default:
			version += fmt.Sprintf(" (%d.%d)", info.dwMajorVersion, info.dwMinorVersion)
		}
	}

	return version, nil
}

func OSversion() string {
	version, err := getWindowsVersion()
	if err != nil {
		fmt.Println("Error:", err)
		return "Error"
	}
	return version
}
