package report

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/cli/cli/v2/pkg/httpmock"
	"github.com/pterm/pterm"
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
		want      []string
	}{{
		name: "empty response body from GitHub API",
		httpStubs: func(reg *httpmock.Registry) {
			createMockRegistry(reg, `{}`)
		},
		want: []string{
			pterm.NewStyle(pterm.FgLightMagenta, pterm.BgBlack, pterm.Bold).Sprintln("OWNER/REPO "),
			"Published <nil>",
			pterm.NewStyle(pterm.FgBlue, pterm.Bold, pterm.Underscore).Sprintln(""),
			"No release assets\n",
			pterm.LightMagenta("0") + " downloads",
		},
	}, {
		name: "when the release has no assets",
		httpStubs: func(reg *httpmock.Registry) {
			createMockRegistry(reg, `{
				"html_url": "https://github.com/FOO/BAR/releases/v1.0.0",
				"tag_name": "v1.0.0",
				"name": "v1.0.0",
				"published_at": "2013-02-27T19:35:32Z",
				"assets": []
			}`)
		},
		want: []string{
			pterm.NewStyle(pterm.FgLightMagenta, pterm.BgBlack, pterm.Bold).Sprintln("OWNER/REPO v1.0.0"),
			"Published 2013-02-27 19:35:32 +0000 UTC",
			pterm.NewStyle(pterm.FgBlue, pterm.Bold, pterm.Underscore).Sprintln("https://github.com/FOO/BAR/releases/v1.0.0"),
			"No release assets\n",
			pterm.LightMagenta("0") + " downloads",
		},
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

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got:\n%q\nwant:\n%q", got, tt.want)
			}

			reg.Verify(t)
		})
	}
}
