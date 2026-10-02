package report

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

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

	tests := []struct {
		name    string
		tag     string
		status  int
		body    string
		wantErr bool
		errMsg  string
		want    []string
	}{{
		name: "empty response body from GitHub API",
		body: `{}`,
		want: []string{
			pterm.NewStyle(pterm.FgLightMagenta, pterm.BgBlack, pterm.Bold).Sprintln("OWNER/REPO "),
			"Published <nil>",
			pterm.NewStyle(pterm.FgBlue, pterm.Bold, pterm.Underscore).Sprintln(""),
			"No release assets\n",
			pterm.LightMagenta("0") + " downloads",
		},
	}, {
		name: "when the release has no assets",
		body: `{
				"html_url": "https://github.com/FOO/BAR/releases/v1.0.0",
				"tag_name": "v1.0.0",
				"name": "v1.0.0",
				"published_at": "2013-02-27T19:35:32Z",
				"assets": []
			}`,
		want: []string{
			pterm.NewStyle(pterm.FgLightMagenta, pterm.BgBlack, pterm.Bold).Sprintln("OWNER/REPO v1.0.0"),
			"Published 2013-02-27 19:35:32 +0000 UTC",
			pterm.NewStyle(pterm.FgBlue, pterm.Bold, pterm.Underscore).Sprintln("https://github.com/FOO/BAR/releases/v1.0.0"),
			"No release assets\n",
			pterm.LightMagenta("0") + " downloads",
		},
	}, {
		name: "when the release has assets",
		body: successfulResp("v1.0.0"),
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
		body: successfulResp("v2.0.0"),
		want: []string{
			pterm.NewStyle(pterm.FgLightMagenta, pterm.BgBlack, pterm.Bold).Sprintln("OWNER/REPO v2.0.0"),
			"Published 2013-02-27 19:35:32 +0000 UTC",
			pterm.NewStyle(pterm.FgBlue, pterm.Bold, pterm.Underscore).Sprintln("https://github.com/FOO/BAR/releases/v2.0.0"),
			barChartStr,
			pterm.LightMagenta("15") + " downloads",
		},
	}, {
		name:    "when the release is not found",
		tag:     "nonexistent",
		status:  http.StatusNotFound,
		body:    `{"message": "Not Found"}`,
		wantErr: true,
		errMsg:  "HTTP 404: Not Found",
		want:    []string{},
	}}

	for _, tt := range tests {
		if tt.tag == "" {
			tt.tag = "latest"
		}
		if tt.status == 0 {
			tt.status = http.StatusOK
		}

		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v3/repos/%s/releases/tags/%s", repo.RepoFullName(), tt.tag)
			if tt.tag == "latest" {
				wantPath = fmt.Sprintf("/api/v3/repos/%s/releases/latest", repo.RepoFullName())
			}

			requests := 0
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, wantPath, r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				fmt.Fprint(w, tt.body)
			}))
			defer server.Close()

			r := *repo
			r.Host = strings.TrimPrefix(server.URL, "https://")

			got, err := Run(&RunOptions{
				Tag:        tt.tag,
				Repo:       &r,
				HTTPClient: server.Client(),
			})

			if tt.wantErr {
				assert.ErrorContains(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got:\n%q\nwant:\n%q", got, tt.want)
			}

			assert.Equal(t, 1, requests)
		})
	}
}

func TestRestURL(t *testing.T) {
	tests := []struct {
		host string
		want string
	}{
		{host: "github.com", want: "https://api.github.com/repos/o/r"},
		{host: "ghes.example.com", want: "https://ghes.example.com/api/v3/repos/o/r"},
		{host: "tenant.ghe.com", want: "https://api.tenant.ghe.com/repos/o/r"},
	}

	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			assert.Equal(t, tt.want, restURL(tt.host, "repos/o/r"))
		})
	}
}
