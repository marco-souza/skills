package cmd

import (
	"fmt"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().Bool("all", false, "Install all skills to target")
	installCmd.Flags().StringP("category", "c", "", "Install all skills matching category")
	installCmd.Flags().Bool("dry-run", false, "Show what would be installed without copying")
}

var installCmd = &cobra.Command{
	Use:     "install [skill|repo] [target]",
	Aliases: []string{"i"},
	Short:   "Install skill(s) to a target project",
	Long: `Install a skill to a target project's .agents/skills directory.

Supports local skills, GitHub shorthand (user/repo), and full GitHub URLs.

Examples:
  skills install git-commit-formatter ./my-project
  skills install --all ./my-project
  skills install --category git ./my-project
  skills install marco/skills ./my-project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		category, _ := cmd.Flags().GetString("category")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		root, _ := cmd.Flags().GetString("root")

		installer := &skills.Installer{DryRun: dryRun, SourceDir: root}

		// When --all or --category, first arg is target
		if all || category != "" {
			target := root
			if len(args) > 0 {
				target = args[0]
			}

			if all {
				return installer.InstallAll(target)
			}

			loader := skills.NewLoader(root)
			allSkills, err := loader.LoadAll()
			if err != nil {
				return err
			}
			filtered := skills.FilterSkills(allSkills, category, nil)
			if len(filtered) == 0 {
				return fmt.Errorf("no skills found in category %q", category)
			}
			for _, s := range filtered {
				if err := installer.Install(s.Name, target); err != nil {
					return err
				}
			}
			return nil
		}

		// Standard install: args[0] = skill, args[1] = target
		if len(args) < 1 {
			return fmt.Errorf("requires at least a skill name or --all/--category flag")
		}

		skillName := args[0]
		target := root
		if len(args) > 1 {
			target = args[1]
		}

		// Try to resolve as a source (GitHub or local)
		src, err := skills.ResolvePath(skillName)
		if err != nil {
			return err
		}

		if src.Type() == "github" {
			ghSrc := src.(*skills.GitHubSource)
			return installer.InstallFromGitHub(ghSrc, target)
		}

		return installer.Install(skillName, target)
	},
}
