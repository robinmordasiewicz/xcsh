package main

import (
	"os"

	"github.com/robinmordasiewicz/vesctl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
