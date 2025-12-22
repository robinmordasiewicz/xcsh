package main

import (
	"os"

	"github.com/robinmordasiewicz/xcsh/cmd"
	"github.com/robinmordasiewicz/xcsh/pkg/errors"
)

func main() {
	if err := cmd.Execute(); err != nil {
		// Use granular exit codes for better scriptability
		os.Exit(errors.GetExitCode(err))
	}
}
