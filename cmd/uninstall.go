package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:     "uninstall <skill>...",
	Aliases: []string{"rm", "remove"},
	Short:   "Uninstall skill(s) from a project",
	Long: `Remove one or more skills from a target project directory.

Skills are removed from .agents/skills/ by default. Use -t to replace
.agents with a custom directory (e.g., -t .opencode removes from .opencode/skills/).

Examples:
  skills uninstall git-commit-formatter
  skills uninstall git-commit-formatter pr-review -t ~/my-project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("requires at least one skill name")
		}

		targetFlag, err := cmd.Flags().GetString("target")
		if err != nil {
			return fmt.Errorf("internal error reading --target flag: %w", err)
		}
		target := targetFlag
		if target == "" {
			target = "."
		}

		// When -t is explicitly set, it replaces .agents as the skills parent directory.
		var skillsDir string
		if targetFlag != "" {
			skillsDir = filepath.Join(target, "skills")
		} else {
			skillsDir = skills.ResolveToSkillsDir(target)
		}

		for _, name := range args {
			skillPath := filepath.Join(skillsDir, name)
			if _, err := os.Stat(skillPath); os.IsNotExist(err) {
				return fmt.Errorf("skill %q not found in %s", name, skillsDir)
			}
			if err := os.RemoveAll(skillPath); err != nil {
				return fmt.Errorf("uninstalling %q: %w", name, err)
			}
			fmt.Printf("Uninstalled skill %q from %s\n", name, skillsDir)
		}
		return nil
	},
}
