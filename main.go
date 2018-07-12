package main

import (
	"os"

	"github.com/krpranay/ciscoampbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
