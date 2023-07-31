package installedsoftware

import (
	"github.com/StackExchange/wmi"
)

type Win32_Product struct {
	Name   string
	Vendor string
}

// This function is slow and can take up to 40 seconds to complete
func WmiSoftwareQuery() []Win32_Product {
	var products []Win32_Product
	err := wmi.Query("SELECT Name, Vendor FROM Win32_Product", &products)
	if err != nil {
		panic(err)
	}

	return products
}
