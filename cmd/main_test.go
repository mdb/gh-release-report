//go:build acceptance
// +build acceptance

package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
)

func TestRootAcceptance(t *testing.T) {
	basicOut := heredoc.Doc(`How many times has a GitHub release been downloaded?

gh release-report reports a release's total download count, as well
as the individual download count for each of its assets.

Usage:
  gh release-report [flags]

Flags:
  -h, --help          help for gh
  -R, --repo string   The targeted repository's full name (default "github.com/mdb/gh-release-report")
  -T, --tag string    The release tag (default "latest")
  -v, --version       version for gh
`)
	tests := []struct {
		args    []string
		wantOut string
		errMsg  string
		wantErr bool
	}{{
		args:    []string{"release-report", "--help"},
		wantOut: basicOut,
	}}

	for _, test := range tests {
		t.Run(fmt.Sprintf("when passed '%s'", strings.Join(test.args, " ")), func(t *testing.T) {
			output, err := exec.Command("gh", test.args...).CombinedOutput()

			if test.wantErr {
				assert.EqualError(t, err, test.errMsg)
			} else {
				assert.NoError(t, err)
			}

			if got := string(output); got != test.wantOut {
				t.Errorf("got stdout:\n%q\nwant:\n%q", got, test.wantOut)
			}
		})
	}
}

func TestReleaseReportAcceptance(t *testing.T) {
	tests := []struct {
		args    []string
		wantOut []string
		errMsg  string
		wantErr bool
	}{{
		args: []string{
			"release-report",
			"--repo=mdb/gh-release-report",
			"--tag=0.0.0",
		},
		wantOut: []string{
			"mdb/gh-release-report 0.0.0",
			"Published 2023-01-23",
			"https://github.com/mdb/gh-release-report/releases/tag/0.0.0",
			"darwin-amd64",
			"darwin-arm64",
			"linux-386",
			"linux-amd64",
			"linux-arm64",
			"netbsd-386",
			"netbsd-amd64",
			"windows-386.exe",
			"windows-amd64.exe",
			"windows-arm64.exe",
			"downloads",
		},
	}, {
		args: []string{
			"release-report",
			"--repo=mdb/gh-release-report",
			"--tag=latest",
		},
		wantOut: []string{
			"mdb/gh-release-report",
			"Published",
			"https://github.com/mdb/gh-release-report/releases/tag/",
			"downloads",
		},
	}, {
		args: []string{
			"release-report",
			"--repo=mdb/gh-release-report",
			"--tag=nonexistent-tag",
		},
		wantOut: []string{
			"HTTP 404: Not Found",
		},
		wantErr: true,
		errMsg:  "exit status 1",
	}}

	for _, test := range tests {
		t.Run(fmt.Sprintf("when passed '%s'", strings.Join(test.args, " ")), func(t *testing.T) {
			output, err := exec.Command("gh", test.args...).CombinedOutput()

			if test.wantErr {
				assert.EqualError(t, err, test.errMsg)
			} else {
				assert.NoError(t, err)
			}

			got := string(output)
			for _, out := range test.wantOut {
				if !strings.Contains(got, out) {
					t.Errorf("expected stdout to include:\n%q\ngot:\n%q", out, got)
				}
			}
		})
	}
}
