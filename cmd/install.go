package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
)

var installCmd = newInstallCmd()

func newInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install <skill>...",
		Aliases: []string{"i"},
		Short:   "Install skill(s) to a target project",
		Long: `Install one or more skills to a target project directory.

Skills are written to .agents/skills/ by default. Use -t to replace
.agents with a custom directory (e.g., -t .opencode writes to .opencode/skills/).

Skills are resolved from the local .agents/skills directory by default.
Use --source to install from a GitHub repo (owner/repo) or local path.
Use --all to install every skill from the source.

Examples:
  skills i git-commit-formatter
  skills i git-commit-formatter pr-review -t ~/my-project
  skills i git-commit-formatter --source marco-souza/skills -t ~/my-project
  skills i --all
  skills i --all --source /path/to/skills-collection -t ~/my-project`,
		RunE: func(cmd *cobra.Command, args []string) error {
			targetFlag, err := cmd.Flags().GetString("target")
			if err != nil {
				return fmt.Errorf("internal error reading --target flag: %w", err)
			}
			allFlag, err := cmd.Flags().GetBool("all")
			if err != nil {
				return fmt.Errorf("internal error reading --all flag: %w", err)
			}

			if allFlag && len(args) > 0 {
				return fmt.Errorf("cannot specify skill names with --all")
			}
			if !allFlag && len(args) == 0 {
				return fmt.Errorf("requires at least one skill name, or use --all")
			}

			source, err := cmd.Flags().GetString("source")
			if err != nil {
				return fmt.Errorf("internal error reading --source flag: %w", err)
			}

			target := targetFlag
			if target == "" {
				target = "."
			}

			sourceDir, cleanup, err := skills.ResolveSourceDir(source, cfg.DefaultSource)
			if err != nil {
				return err
			}
			if cleanup != nil {
				defer cleanup()
			}

			installer := &skills.Installer{SourceDir: sourceDir}

			// Helper to determine parent dir
			parentDir := filepath.Join(target, ".agents")

			if allFlag {
				return installer.InstallAll(parentDir)
			}

			for _, name := range args {
				if err := installer.Install(name, parentDir); err != nil {
					return fmt.Errorf("installing %q: %w", name, err)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringP("target", "t", "", "Target project directory")
	cmd.Flags().Bool("all", false, "Install all skills from the source")
	cmd.Flags().StringP("source", "s", "", "Source for skills: GitHub repo (owner/repo) or local path")
	return cmd
}
