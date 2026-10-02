package report

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func NewCmdRoot(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gh release-report",
		Short: "How many times has a GitHub release been downloaded?",
		Long: heredoc.Doc(`
			How many times has a GitHub release been downloaded?

			gh release-report reports a release's total download count, as well
			as the individual download count for each of its assets.
		`),
		SilenceUsage: true,
		Version:      version,
		RunE: func(cmd *cobra.Command, args []string) error {
			tag, err := cmd.Flags().GetString("tag")
			if err != nil {
				return err
			}

			repo, err := getRepoOption(cmd)
			if err != nil {
				return err
			}

			httpClient, err := ghapi.NewHTTPClient(ghapi.ClientOptions{Host: repo.RepoHost()})
			if err != nil {
				return err
			}

			contents, err := Run(&RunOptions{
				Repo:       repo,
				Tag:        tag,
				HTTPClient: httpClient,
			})
			if err != nil {
				return err
			}

			pterm.DefaultBox.Println(strings.Join(contents, "\n"))

			return nil
		},
	}

	defaultRepo := ""
	currentRepo, err := repository.Current()
	if err == nil {
		defaultRepo = fmt.Sprintf("%s/%s/%s", currentRepo.Host, currentRepo.Owner, currentRepo.Name)
	}

	var repo string
	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "R", defaultRepo, "The targeted repository's full name")

	var tag string
	rootCmd.PersistentFlags().StringVarP(&tag, "tag", "T", "latest", "The release tag")

	return rootCmd
}

type RunOptions struct {
	Repo       *ghRepo
	Tag        string
	HTTPClient *http.Client
}

type release struct {
	TagName     string         `json:"tag_name"`
	URL         string         `json:"html_url"`
	PublishedAt *time.Time     `json:"published_at"`
	Assets      []releaseAsset `json:"assets"`
}

type releaseAsset struct {
	Name          string `json:"name"`
	DownloadCount int    `json:"download_count"`
}

func Run(opts *RunOptions) ([]string, error) {
	repo := opts.Repo
	var url string

	switch tag := opts.Tag; tag {
	case "latest":
		url = fmt.Sprintf("repos/%s/releases/latest", repo.RepoFullName())
	default:
		url = fmt.Sprintf("repos/%s/releases/tags/%s", repo.RepoFullName(), tag)
	}

	response, err := fetchRelease(opts.HTTPClient, restURL(repo.RepoHost(), url))
	if err != nil {
		return []string{}, err
	}

	total := 0
	bars := pterm.Bars{}
	for _, asset := range response.Assets {
		// TODO: make --exclude a configurable option
		if strings.Contains(strings.ToLower(asset.Name), "checksums") || strings.Contains(strings.ToLower(asset.Name), "sha256sums") {
			continue
		}

		total += asset.DownloadCount
		bars = append(bars, pterm.Bar{
			Label: asset.Name,
			Value: asset.DownloadCount,
		})
	}

	chart, err := pterm.DefaultBarChart.WithHorizontal().WithBars(bars).WithShowValue().Srender()
	if err != nil {
		return []string{}, err
	}

	if len(response.Assets) == 0 {
		chart = "No release assets\n"
	}

	title := fmt.Sprintf("%s %s", repo.RepoFullName(), response.TagName)
	p := message.NewPrinter(language.English)
	formattedTotal := p.Sprintf("%d", total)
	emphasized := pterm.NewStyle(pterm.FgLightMagenta, pterm.BgBlack, pterm.Bold)

	return []string{
		emphasized.Sprintln(title),
		fmt.Sprintf("Published %s", response.PublishedAt),
		pterm.NewStyle(pterm.FgBlue, pterm.Bold, pterm.Underscore).Sprintln(response.URL),
		chart,
		pterm.LightMagenta(formattedTotal) + " downloads",
	}, nil
}

func fetchRelease(client *http.Client, url string) (*release, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, ghapi.HandleHTTPError(resp)
	}

	var r release
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}

// restURL returns the REST API URL for path on host. GitHub Enterprise Server
// serves the API under /api/v3/; github.com and GHE.com use an api. subdomain.
func restURL(host, path string) string {
	host = auth.NormalizeHostname(host)
	if auth.IsEnterprise(host) {
		return fmt.Sprintf("https://%s/api/v3/%s", host, path)
	}
	return fmt.Sprintf("https://api.%s/%s", host, path)
}
