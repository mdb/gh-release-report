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
		name: "when the release has no assets or known publication date",
		httpStubs: func(reg *httpmock.Registry) {
			createMockRegistry(reg, `{}`)
		},
		wantOut: "\x1b[95;40;1m\x1b[95;40;1mOWNER/REPO \x1b[0m\x1b[0m\n\nPublished <nil>\n\x1b[34;1;4m\x1b[34;1;4m\x1b[0m\x1b[0m\n\nNo release assets\n\n\x1b[95m0\x1b[0m downloads",
	}, {
		name: "when the release has assets",
		httpStubs: func(reg *httpmock.Registry) {
			createMockRegistry(reg, `{
				"html_url": "https://github.com/FOO/BAR/releases/v1.0.0",
				"tag_name": "v1.0.0",
				"name": "v1.0.0",
				"published_at": "2013-02-27T19:35:32Z",
				"assets": [
					{
						"name": "example.zip",
						"download_count": 42
					},
					{
						"name": "example.tar.gz",
						"download_count": 100
					}
				]
			}`)
		},
		wantOut: "\x1b[95;40;1m\x1b[95;40;1mOWNER/REPO v1.0.0\x1b[0m\x1b[0m\n\nPublished 2013-02-27 19:35:32 +0000 UTC\n\x1b[34;1;4m\x1b[34;1;4mhttps://github.com/FOO/BAR/releases/v1.0.0\x1b[0m\x1b[0m\n\n\x1b[0m\x1b[0m               \x1b[0m\x1b[0m                                                           \n\x1b[96m\x1b[96mexample.zip\x1b[0m\x1b[0m\x1b[0m\x1b[0m    \x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m                                 42\x1b[0m\x1b[0m  \n\x1b[96m\x1b[96mexample.tar.gz\x1b[0m\x1b[0m \x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m   100 \n\n\x1b[95m142\x1b[0m downloads",
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
