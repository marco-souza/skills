package cmd

import (
	"fmt"
	"os"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().Bool("all", false, "Install all skills to target")
	installCmd.Flags().StringP("category", "c", "", "Install all skills matching category")
	installCmd.Flags().Bool("dry-run", false, "Show what would be installed without copying")
	installCmd.Flags().StringP("target", "t", "", "Target project directory")
}

var installCmd = &cobra.Command{
	Use:     "install <skill>... [target]",
	Aliases: []string{"i"},
	Short:   "Install skill(s) to a target project",
	Long: `Install one or more skills to a target project's .agents/skills directory.

Supports local skills, GitHub shorthand (user/repo), and full GitHub URLs.

Examples:
  skills i git-commit-formatter pr-review -t ./my-project
  skills install git-commit-formatter ./my-project
  skills i --all -t ./my-project
  skills i --category git -t ./my-project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		category, _ := cmd.Flags().GetString("category")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		targetFlag, _ := cmd.Flags().GetString("target")
		repo, _ := cmd.Flags().GetString("repo")
		root, _ := cmd.Flags().GetString("root")

		// Determine source directory
		sourceDir := root

		// Prefer local skills if they exist
		localSkills := skills.ResolveToSkillsDir(root)
		hasLocal := false
		if _, err := os.Stat(localSkills); err == nil {
			hasLocal = true
		}

		// Use --repo only when no local skills exist
		if !hasLocal && repo != "" {
			src, err := skills.ResolvePath(repo)
			if err != nil {
				return fmt.Errorf("resolving repo: %w", err)
			}
			if src.Type() != "github" {
				return fmt.Errorf("--repo requires a GitHub source, got %s", src.Type())
			}
			repoDir, cleanup, err := skills.CloneRepo(src.(*skills.GitHubSource))
			if err != nil {
				return err
			}
			defer cleanup()
			sourceDir = repoDir
		}

		installer := &skills.Installer{DryRun: dryRun, SourceDir: sourceDir}

		// --all or --category: install bulk
		if all || category != "" {
			target := targetFlag
			if target == "" {
				target = root
				if len(args) > 0 {
					target = args[0]
				}
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

		// No skills specified
		if len(args) == 0 {
			return fmt.Errorf("requires at least a skill name, or use --all/--category")
		}

		// Determine target: -t flag > last arg (if it looks like a path) > root
		target := targetFlag
		var skillNames []string

		if target != "" {
			// Target explicitly set via flag, all args are skills
			skillNames = args
		} else if len(args) == 1 {
			// Single arg: could be skill name or GitHub shorthand — treat as skill
			skillNames = args
			target = root
		} else {
			// Multiple args: last arg is target, rest are skills
			skillNames = args[:len(args)-1]
			target = args[len(args)-1]
		}

		// Install each skill
		for _, name := range skillNames {
			src, err := skills.ResolvePath(name)
			if err != nil {
				return err
			}

			if src.Type() == "github" {
				if err := installer.InstallFromGitHub(src.(*skills.GitHubSource), target); err != nil {
					return err
				}
				continue
			}

			if err := installer.Install(name, target); err != nil {
				return err
			}
		}

		return nil
	},
}
