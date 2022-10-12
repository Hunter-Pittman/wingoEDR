package main

import (
	"fmt"
	"wingoEDR/processes"
)

func main() {
	fmt.Printf("%v", processes.GetAllProcesses())
}
