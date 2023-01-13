package report

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	shared "github.com/cli/cli/v2/pkg/cmd/release/shared"
	gh "github.com/cli/go-gh"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewCmdRoot(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gh release-report",
		Short: "TODO",
		Long: heredoc.Doc(`
			TODO
		`),
		SilenceUsage: true,
		Version:      version,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := getRepoOption(cmd)
			if err != nil {
				return err
			}

			client, err := gh.RESTClient(nil)
			if err != nil {
				return err
			}

			var response shared.Release
			err = client.Get(fmt.Sprintf("repos/%s/%s/releases/latest", repo.Owner, repo.Name), &response)
			if err != nil {
				return err
			}

			total := 0
			bars := pterm.Bars{}
			for _, asset := range response.Assets {
				if strings.Contains(asset.Name, "checksums") {
					continue
				}

				total += asset.DownloadCount
				bars = append(bars, pterm.Bar{
					Label: asset.Name,
					Value: asset.DownloadCount,
				})
			}

			pterm.DefaultBasicText.Println(fmt.Sprintf("%d total downloads", total))
			return pterm.DefaultBarChart.WithHorizontal().WithBars(bars).WithShowValue().Render()
		},
	}

	defaultRepo := ""
	currentRepo, _ := gh.CurrentRepository()
	if currentRepo != nil {
		defaultRepo = fmt.Sprintf("%s/%s/%s", currentRepo.Host(), currentRepo.Owner(), currentRepo.Name())
	}

	var repo string
	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "R", defaultRepo, "The targeted repository's full name")

	return rootCmd
}
