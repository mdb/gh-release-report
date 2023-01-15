package report

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
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
			tag, _ := cmd.Flags().GetString("tag")

			repo, err := getRepoOption(cmd)
			if err != nil {
				return err
			}

			url := fmt.Sprintf("repos/%s/releases/latest", repo.RepoFullName())
			if tag != "latest" {
				url = fmt.Sprintf("repos/%s/releases/tags/%s", repo.RepoFullName(), tag)
			}

			client, err := gh.RESTClient(nil)
			if err != nil {
				return err
			}

			var response shared.Release
			err = client.Get(url, &response)
			if err != nil {
				return err
			}

			total := 0
			bars := pterm.Bars{}
			for _, asset := range response.Assets {
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
				return err
			}

			title := fmt.Sprintf("%s %s", repo.RepoFullName(), response.TagName)
			p := message.NewPrinter(language.English)
			totalDs := p.Sprintf("%d downloads", total)
			published := fmt.Sprintf("Published %s", response.PublishedAt)

			pterm.DefaultBox.WithTitle(title).Println("\n" + published + "\n" + response.URL + "\n" + chart + "\n" + totalDs)

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
