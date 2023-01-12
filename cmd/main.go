package main

import (
	"os"

	"github.com/mdb/gh-release-report/internal/report"
)

// version's value is passed in at build time.
var version string

func main() {
	rootCmd := report.NewCmdRoot(version)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
