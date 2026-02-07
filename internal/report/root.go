package report

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	tea "github.com/charmbracelet/bubbletea"
	cliapi "github.com/cli/cli/v2/api"
	shared "github.com/cli/cli/v2/pkg/cmd/release/shared"
	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
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

			ghClient, err := ghapi.DefaultHTTPClient()
			if err != nil {
				return err
			}

			data, err := FetchRelease(&RunOptions{
				Repo:       repo,
				Tag:        tag,
				HTTPClient: ghClient,
			})
			if err != nil {
				return err
			}

			m := newModel(data)
			p := tea.NewProgram(m, tea.WithOutput(os.Stdout), tea.WithInput(nil))
			_, err = p.Run()
			return err
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

// FetchRelease retrieves release data from the GitHub API.
func FetchRelease(opts *RunOptions) (*ReleaseData, error) {
	repo := opts.Repo
	var url string

	switch tag := opts.Tag; tag {
	case "latest":
		url = fmt.Sprintf("repos/%s/releases/latest", repo.RepoFullName())
	default:
		url = fmt.Sprintf("repos/%s/releases/tags/%s", repo.RepoFullName(), tag)
	}

	ghClient := cliapi.NewClientFromHTTP(opts.HTTPClient)
	var response shared.Release
	err := ghClient.REST(repo.RepoHost(), "GET", url, nil, &response)
	if err != nil {
		return nil, err
	}

	data := &ReleaseData{
		RepoFullName: repo.RepoFullName(),
		TagName:      response.TagName,
		PublishedAt:  response.PublishedAt,
		URL:          response.URL,
		Assets:       []AssetData{},
		TotalCount:   0,
	}

	for _, asset := range response.Assets {
		// TODO: make --exclude a configurable option
		if strings.Contains(strings.ToLower(asset.Name), "checksums") || strings.Contains(strings.ToLower(asset.Name), "sha256sums") {
			continue
		}

		data.TotalCount += asset.DownloadCount
		data.Assets = append(data.Assets, AssetData{
			Name:          asset.Name,
			DownloadCount: asset.DownloadCount,
		})
	}

	return data, nil
}
