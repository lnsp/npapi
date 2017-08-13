package main

import (
	"fmt"

	"github.com/lnsp/npapi"
)

const (
	addr = "0x6fc554f5a9a7b4ce02922ed98e80b804c2101c4b"
)

func main() {
	balance, err := npapi.Balance(addr)
	fmt.Println(balance, err)

	averageIn, err := npapi.AverageHashrateIn(addr, 3)
	fmt.Println(averageIn, err)

	average, err := npapi.AverageHashrate(addr)
	fmt.Println(average, err)

	chart, err := npapi.HashrateChart(addr)
	fmt.Println(chart, err)

	exists := npapi.Exists(addr)
	fmt.Println(exists)

	current, err := npapi.CurrentHashrate(addr)
	fmt.Println(current, err)

	history, err := npapi.HashrateHistory(addr)
	fmt.Println(history, err)

	user, err := npapi.UserInfo(addr)
	fmt.Println(user, err)
}
