# npapi [![GoDoc](https://godoc.org/github.com/lnsp/npapi?status.svg)](https://godoc.org/github.com/lnsp/npapi)

A lightweight Go wrapper for the Nanopool Ethereum API.

## Example
````go
package main

import (
	"fmt"
	"github.com/lnsp/npapi"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "USAGE: go run main.go [nanopool account address]")
		return
	}
	balance, err := npapi.Balance(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch balance: %v\n", err)
		return
	}
	fmt.Printf("You have earned %.6f ETH.\n", balance)
}
```
