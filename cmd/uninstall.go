package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().StringP("target", "t", "", "Target project directory")
}

var uninstallCmd = &cobra.Command{
	Use:     "uninstall <skill>...",
	Aliases: []string{"rm", "remove"},
	Short:   "Uninstall skill(s) from a project",
	Long: `Remove one or more skills from a project's .agents/skills directory.

Examples:
  skills uninstall git-commit-formatter
  skills uninstall git-commit-formatter pr-review -t ~/my-project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("requires at least one skill name")
		}

		targetFlag, _ := cmd.Flags().GetString("target")
		target := targetFlag
		if target == "" {
			target = "."
		}

		skillsDir := skills.ResolveToSkillsDir(target)

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
