package report

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/cli/cli/v2/pkg/httpmock"
	"github.com/pterm/pterm"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	repo := &ghRepo{
		Owner: "OWNER",
		Name:  "REPO",
	}

	successfulResp := func(tag string) string {
		return fmt.Sprintf(`{
			"html_url": "https://github.com/FOO/BAR/releases/%s",
			"tag_name": "%s",
			"name": "%s",
			"published_at": "2013-02-27T19:35:32Z",
			"assets": [{
				"name": "example.zip",
				"download_count": 10
			}, {
				"name": "exampletwo.zip",
				"download_count": 5
			}]
		}`, tag, tag, tag)
	}

	barChartStr := "\x1b[0m\x1b[0m               \x1b[0m\x1b[0m                                                          \n\x1b[96m\x1b[96mexample.zip\x1b[0m\x1b[0m\x1b[0m\x1b[0m    \x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m  10\x1b[0m\x1b[0m \n\x1b[96m\x1b[96mexampletwo.zip\x1b[0m\x1b[0m \x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m\x1b[36m\x1b[36m█\x1b[0m\x1b[0m                             5  \n"

	createMockRegistry := func(reg *httpmock.Registry, tag, resp string) {
		var url string
		switch tag {
		case "latest":
			url = fmt.Sprintf("repos/%s/releases/latest", repo.RepoFullName())
		default:
			url = fmt.Sprintf("repos/%s/releases/tags/%s", repo.RepoFullName(), tag)
		}

		reg.Register(
			httpmock.REST("GET", url),
			httpmock.StringResponse(resp))
	}

	tests := []struct {
		name      string
		tag       string
		httpStubs func(*httpmock.Registry)
		wantErr   bool
		errMsg    string
		want      []string
	}{{
		name: "empty response body from GitHub API",
		httpStubs: func(reg *httpmock.Registry) {
			createMockRegistry(reg, "latest", `{}`)
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
			createMockRegistry(reg, "latest", `{
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
	}, {
		name: "when the release has assets",
		httpStubs: func(reg *httpmock.Registry) {
			createMockRegistry(reg, "latest", successfulResp("v1.0.0"))
		},
		want: []string{
			pterm.NewStyle(pterm.FgLightMagenta, pterm.BgBlack, pterm.Bold).Sprintln("OWNER/REPO v1.0.0"),
			"Published 2013-02-27 19:35:32 +0000 UTC",
			pterm.NewStyle(pterm.FgBlue, pterm.Bold, pterm.Underscore).Sprintln("https://github.com/FOO/BAR/releases/v1.0.0"),
			barChartStr,
			pterm.LightMagenta("15") + " downloads",
		},
	}, {
		name: "when a tag is specified",
		tag:  "v2.0.0",
		httpStubs: func(reg *httpmock.Registry) {
			createMockRegistry(reg, "v2.0.0", successfulResp("v2.0.0"))
		},
		want: []string{
			pterm.NewStyle(pterm.FgLightMagenta, pterm.BgBlack, pterm.Bold).Sprintln("OWNER/REPO v2.0.0"),
			"Published 2013-02-27 19:35:32 +0000 UTC",
			pterm.NewStyle(pterm.FgBlue, pterm.Bold, pterm.Underscore).Sprintln("https://github.com/FOO/BAR/releases/v2.0.0"),
			barChartStr,
			pterm.LightMagenta("15") + " downloads",
		},
	}}

	for _, tt := range tests {
		reg := &httpmock.Registry{}
		tt.httpStubs(reg)
		if tt.tag == "" {
			tt.tag = "latest"
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := Run(&RunOptions{
				Tag:  tt.tag,
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
