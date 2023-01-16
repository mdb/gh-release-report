package report

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc"
	cliapi "github.com/cli/cli/v2/api"
	shared "github.com/cli/cli/v2/pkg/cmd/release/shared"
	gh "github.com/cli/go-gh"
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

			ghClient, _ := cliapi.NewHTTPClient(cliapi.HTTPClientOptions{
				Config: &Config{},
			})

			contents, err := Run(&RunOptions{
				Repo:       repo,
				Tag:        tag,
				HTTPClient: ghClient,
			})
			if err != nil {
				return err
			}

			pterm.DefaultBox.Println(contents)

			return nil
		},
	}

	defaultRepo := ""
	currentRepo, _ := gh.CurrentRepository()
	if currentRepo != nil {
		defaultRepo = fmt.Sprintf("%s/%s/%s", currentRepo.Host(), currentRepo.Owner(), currentRepo.Name())
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

func Run(opts *RunOptions) (string, error) {
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
		return "", err
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
		return "", err
	}

	if len(response.Assets) == 0 {
		chart = "No release assets\n"
	}

	title := fmt.Sprintf("%s %s", repo.RepoFullName(), response.TagName)
	p := message.NewPrinter(language.English)
	formattedTotal := p.Sprintf("%d", total)
	emphasized := pterm.NewStyle(pterm.FgLightMagenta, pterm.BgBlack, pterm.Bold)

	return strings.Join([]string{
		emphasized.Sprintln(title),
		fmt.Sprintf("Published %s", response.PublishedAt),
		pterm.NewStyle(pterm.FgBlue, pterm.Bold, pterm.Underscore).Sprintln(response.URL),
		chart,
		pterm.LightMagenta(formattedTotal) + " downloads",
	}, "\n"), nil
}
