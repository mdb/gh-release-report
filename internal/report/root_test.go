package report

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cli/cli/v2/pkg/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestFetchRelease(t *testing.T) {
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

	publishedAt := time.Date(2013, 2, 27, 19, 35, 32, 0, time.UTC)

	tests := []struct {
		name      string
		tag       string
		httpStubs func(*httpmock.Registry)
		wantErr   bool
		errMsg    string
		want      *ReleaseData
	}{{
		name: "empty response body from GitHub API",
		httpStubs: func(reg *httpmock.Registry) {
			createMockRegistry(reg, "latest", `{}`)
		},
		want: &ReleaseData{
			RepoFullName: "OWNER/REPO",
			TagName:      "",
			PublishedAt:  nil,
			URL:          "",
			Assets:       []AssetData{},
			TotalCount:   0,
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
		want: &ReleaseData{
			RepoFullName: "OWNER/REPO",
			TagName:      "v1.0.0",
			PublishedAt:  &publishedAt,
			URL:          "https://github.com/FOO/BAR/releases/v1.0.0",
			Assets:       []AssetData{},
			TotalCount:   0,
		},
	}, {
		name: "when the release has assets",
		httpStubs: func(reg *httpmock.Registry) {
			createMockRegistry(reg, "latest", successfulResp("v1.0.0"))
		},
		want: &ReleaseData{
			RepoFullName: "OWNER/REPO",
			TagName:      "v1.0.0",
			PublishedAt:  &publishedAt,
			URL:          "https://github.com/FOO/BAR/releases/v1.0.0",
			Assets: []AssetData{
				{Name: "example.zip", DownloadCount: 10},
				{Name: "exampletwo.zip", DownloadCount: 5},
			},
			TotalCount: 15,
		},
	}, {
		name: "when a tag is specified",
		tag:  "v2.0.0",
		httpStubs: func(reg *httpmock.Registry) {
			createMockRegistry(reg, "v2.0.0", successfulResp("v2.0.0"))
		},
		want: &ReleaseData{
			RepoFullName: "OWNER/REPO",
			TagName:      "v2.0.0",
			PublishedAt:  &publishedAt,
			URL:          "https://github.com/FOO/BAR/releases/v2.0.0",
			Assets: []AssetData{
				{Name: "example.zip", DownloadCount: 10},
				{Name: "exampletwo.zip", DownloadCount: 5},
			},
			TotalCount: 15,
		},
	}, {
		name: "checksums assets are excluded",
		httpStubs: func(reg *httpmock.Registry) {
			createMockRegistry(reg, "latest", `{
				"html_url": "https://github.com/FOO/BAR/releases/v1.0.0",
				"tag_name": "v1.0.0",
				"name": "v1.0.0",
				"published_at": "2013-02-27T19:35:32Z",
				"assets": [{
					"name": "example.zip",
					"download_count": 10
				}, {
					"name": "checksums.txt",
					"download_count": 100
				}, {
					"name": "SHA256SUMS",
					"download_count": 50
				}]
			}`)
		},
		want: &ReleaseData{
			RepoFullName: "OWNER/REPO",
			TagName:      "v1.0.0",
			PublishedAt:  &publishedAt,
			URL:          "https://github.com/FOO/BAR/releases/v1.0.0",
			Assets: []AssetData{
				{Name: "example.zip", DownloadCount: 10},
			},
			TotalCount: 10,
		},
	}}

	for _, tt := range tests {
		reg := &httpmock.Registry{}
		tt.httpStubs(reg)
		if tt.tag == "" {
			tt.tag = "latest"
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchRelease(&RunOptions{
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

			assert.Equal(t, tt.want, got)

			reg.Verify(t)
		})
	}
}

func TestRenderBarChart(t *testing.T) {
	tests := []struct {
		name         string
		assets       []AssetData
		wantEmpty    bool
		wantContains []string
	}{{
		name:      "empty assets",
		assets:    []AssetData{},
		wantEmpty: true,
	}, {
		name: "single asset",
		assets: []AssetData{
			{Name: "file.zip", DownloadCount: 100},
		},
		wantContains: []string{"file.zip", "██████████████████████████████████████████████████", "100"},
	}, {
		name: "multiple assets",
		assets: []AssetData{
			{Name: "file.zip", DownloadCount: 100},
			{Name: "other.tar.gz", DownloadCount: 50},
		},
		wantContains: []string{
			"file.zip",
			"other.tar.gz",
			"██████████████████████████████████████████████████", // 50 chars for max
			"█████████████████████████",                          // 25 chars for half
			"100",
			"50",
		},
	}, {
		name: "zero count shows no bar",
		assets: []AssetData{
			{Name: "file.zip", DownloadCount: 100},
			{Name: "empty.zip", DownloadCount: 0},
		},
		wantContains: []string{"file.zip", "empty.zip", "100", " 0"},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderBarChart(tt.assets)
			if tt.wantEmpty {
				assert.Empty(t, got)
			} else {
				for _, s := range tt.wantContains {
					assert.Contains(t, got, s)
				}
			}
		})
	}
}
