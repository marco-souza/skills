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
	Long:    `List skills from the local .agents/skills directory, a local path, or a remote GitHub repo.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root := "."
		if len(args) > 0 {
			root = args[0]
		}

		source, _ := cmd.Flags().GetString("source")

		// Try local first (no --source flag, look in .agents/skills)
		if source == "" {
			localDir := skills.ResolveToSkillsDir(root)
			if _, err := os.Stat(localDir); err == nil {
				return listSkills(root)
			}
			source = cfg.DefaultSource
		}

		if source == "" {
			localDir := skills.ResolveToSkillsDir(root)
			return fmt.Errorf("no local skills found at %s and no --source specified", localDir)
		}

		src, err := skills.ResolvePath(source)
		if err != nil {
			return fmt.Errorf("resolving source %q: %w", source, err)
		}

		switch s := src.(type) {
		case *skills.LocalSource:
			return listSkills(s.Path)
		case *skills.GitHubSource:
			tmpDir, cleanup, err := skills.CloneRepo(s)
			if err != nil {
				return err
			}
			defer cleanup()
			return listSkills(tmpDir)
		default:
			return fmt.Errorf("--source requires a GitHub repo (owner/repo) or local path, got %s", source)
		}
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
