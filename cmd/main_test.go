//go:build acceptance
// +build acceptance

package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	// compile a 'gh-release-report' for for use in running tests
	exe := exec.Command("go", "build", "-ldflags", "-X main.version=test", "-o", "gh-release-report")
	err := exe.Run()
	if err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestAcceptance(t *testing.T) {

}
