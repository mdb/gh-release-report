package report

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	gh "github.com/cli/go-gh"
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
			r, _ := cmd.Flags().GetString("repo")
			fmt.Println(r)
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

	return rootCmd
}
