package cmd

import (
	"fmt"
	"os"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringP("target", "t", "", "Target project directory")
}

var installCmd = &cobra.Command{
	Use:     "install <skill>...",
	Aliases: []string{"i"},
	Short:   "Install skill(s) to a target project",
	Long: `Install one or more skills to a target project's .agents/skills directory.

Skills are resolved from the local .agents/skills directory by default.
Use --repo to install from a remote GitHub repository.

Examples:
  skills i git-commit-formatter
  skills i git-commit-formatter pr-review -t ~/my-project
  skills i git-commit-formatter -r marco-souza/skills -t ~/my-project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("requires at least one skill name")
		}

		targetFlag, _ := cmd.Flags().GetString("target")
		repo, _ := cmd.Flags().GetString("repo")
		if repo == "" {
			repo = cfg.DefaultRepo
		}

		target := targetFlag
		if target == "" {
			target = "."
		}

		sourceDir, cleanup, err := resolveSource(repo)
		if err != nil {
			return err
		}
		if cleanup != nil {
			defer cleanup()
		}

		installer := &skills.Installer{SourceDir: sourceDir}

		for _, name := range args {
			if err := installer.Install(name, target); err != nil {
				return fmt.Errorf("installing %q: %w", name, err)
			}
		}
		return nil
	},
}

// resolveSource returns the directory to install skills from and an optional cleanup function.
// Prefers local .agents/skills, falls back to cloning the remote repo.
func resolveSource(repo string) (string, func(), error) {
	localDir := skills.ResolveToSkillsDir(".")
	if _, err := os.Stat(localDir); err == nil {
		return ".", nil, nil
	}

	if repo == "" {
		return "", nil, fmt.Errorf("no local skills found and no --repo specified")
	}

	src, err := skills.ResolvePath(repo)
	if err != nil {
		return "", nil, fmt.Errorf("resolving repo %q: %w", repo, err)
	}
	gh, ok := src.(*skills.GitHubSource)
	if !ok {
		return "", nil, fmt.Errorf("--repo requires a GitHub repo (owner/repo), got %s", repo)
	}

	tmpDir, cleanup, err := skills.CloneRepo(gh)
	if err != nil {
		return "", nil, err
	}
	return tmpDir, cleanup, nil
}
