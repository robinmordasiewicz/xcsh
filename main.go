package main

import (
	"os"

	"github.com/robinmordasiewicz/f5xc/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
