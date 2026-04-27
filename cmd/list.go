package cmd

import (
	"fmt"
	"os"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list [path]",
	Aliases: []string{"ls"},
	Short:   "List available skills",
	Long:    `List skills from the local .agents/skills directory or a remote GitHub repo.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root := "."
		if len(args) > 0 {
			root = args[0]
		}

		repo, _ := cmd.Flags().GetString("repo")
		if repo == "" {
			repo = cfg.DefaultRepo
		}

		// Try local first
		localDir := skills.ResolveToSkillsDir(root)
		if _, err := os.Stat(localDir); err == nil {
			return listSkills(root)
		}

		// Fall back to remote repo
		if repo == "" {
			return fmt.Errorf("no local skills found at %s and no --repo specified", localDir)
		}

		src, err := skills.ResolvePath(repo)
		if err != nil {
			return fmt.Errorf("resolving repo %q: %w", repo, err)
		}
		gh, ok := src.(*skills.GitHubSource)
		if !ok {
			return fmt.Errorf("--repo requires a GitHub repo (owner/repo), got %s", repo)
		}

		tmpDir, cleanup, err := skills.CloneRepo(gh)
		if err != nil {
			return err
		}
		defer cleanup()

		return listSkills(tmpDir)
	},
}

func listSkills(root string) error {
	loader := skills.NewLoader(root)
	sk, err := loader.LoadAll()
	if err != nil {
		return err
	}

	if len(sk) == 0 {
		fmt.Println("no skills found")
		return nil
	}

	for _, s := range sk {
		fmt.Println(s.Name)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
