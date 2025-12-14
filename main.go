package main

import (
	"os"

	"github.com/robinmordasiewicz/f5xcctl/cmd"
	"github.com/robinmordasiewicz/f5xcctl/pkg/errors"
)

func main() {
	if err := cmd.Execute(); err != nil {
		// Use granular exit codes for better scriptability
		os.Exit(errors.GetExitCode(err))
	}
}
