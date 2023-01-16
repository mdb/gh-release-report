package report

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/cli/cli/v2/pkg/httpmock"
	"github.com/stretchr/testify/assert"
)

func execute(t *testing.T, args string) string {
	t.Helper()

	cmd := NewCmdRoot("test")

	actual := new(bytes.Buffer)

	cmd.SetOut(actual)
	cmd.SetErr(actual)
	cmd.SetArgs(strings.Split(args, " "))
	cmd.Execute()

	return actual.String()
}

func TestRun(t *testing.T) {
	repo := &ghRepo{
		Owner: "OWNER",
		Name:  "REPO",
	}

	createMockRegistry := func(reg *httpmock.Registry, resp string) {
		reg.Register(
			httpmock.REST("GET", fmt.Sprintf("repos/%s/releases/latest", repo.RepoFullName())),
			httpmock.StringResponse(resp))
	}

	tests := []struct {
		name      string
		httpStubs func(*httpmock.Registry)
		wantErr   bool
		errMsg    string
		wantOut   string
	}{{
		name: "when the release has no assets",
		httpStubs: func(reg *httpmock.Registry) {
			createMockRegistry(reg, `{}`)
		},
		wantOut: "\x1b[95;40;1m\x1b[95;40;1mOWNER/REPO \x1b[0m\x1b[0m\n\nPublished <nil>\n\x1b[34;1;4m\x1b[34;1;4m\x1b[0m\x1b[0m\n\nNo release assets\n\n\x1b[95m0\x1b[0m downloads",
	}}

	for _, tt := range tests {
		reg := &httpmock.Registry{}
		tt.httpStubs(reg)

		t.Run(tt.name, func(t *testing.T) {
			got, err := Run(&RunOptions{
				Tag:  "latest",
				Repo: repo,
				HTTPClient: &http.Client{
					Transport: reg,
				},
			})

			if tt.wantErr {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			if got != tt.wantOut {
				t.Errorf("got:\n%q\nwant:\n%q", got, tt.wantOut)
			}

			reg.Verify(t)
		})
	}
}
