package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringP("target", "t", "", "Target project directory")
	installCmd.Flags().Bool("all", false, "Install all skills from the source")
}

var installCmd = &cobra.Command{
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
		targetFlag, _ := cmd.Flags().GetString("target")
		allFlag, _ := cmd.Flags().GetBool("all")

		if allFlag && len(args) > 0 {
			return fmt.Errorf("cannot specify skill names with --all")
		}
		if !allFlag && len(args) == 0 {
			return fmt.Errorf("requires at least one skill name, or use --all")
		}

		source, _ := cmd.Flags().GetString("source")

		target := targetFlag
		if target == "" {
			target = "."
		}

		sourceDir, cleanup, err := resolveSource(source)
		if err != nil {
			return err
		}
		if cleanup != nil {
			defer cleanup()
		}

		installer := &skills.Installer{SourceDir: sourceDir}

		// Helper to determine parent dir
		parentDir := func() string {
			if targetFlag != "" {
				return target
			}
			return filepath.Join(target, ".agents")
		}()

		if allFlag {
			return installAll(installer, sourceDir, parentDir)
		}

		for _, name := range args {
			if err := installer.Install(name, parentDir); err != nil {
				return fmt.Errorf("installing %q: %w", name, err)
			}
		}
		return nil
	},
}

// resolveSource returns the directory containing skills/ and an optional cleanup function.
// First checks local .agents/skills, then falls back to --source (local path or GitHub repo).
func resolveSource(source string) (string, func(), error) {
	// No explicit --source: prefer local .agents/skills, fall back to default repo
	if source == "" {
		localDir := skills.ResolveToSkillsDir(".")
		if _, err := os.Stat(localDir); err == nil {
			return ".", nil, nil
		}
		// No local skills, use default source (GitHub repo by default)
		source = cfg.DefaultSource
		if source == "" {
			return "", nil, fmt.Errorf("no local skills found and no --source specified")
		}
	}

	// Resolve the source: local path or GitHub repo
	src, err := skills.ResolvePath(source)
	if err != nil {
		return "", nil, fmt.Errorf("resolving source %q: %w", source, err)
	}

	switch s := src.(type) {
	case *skills.LocalSource:
		return s.Path, nil, nil
	case *skills.GitHubSource:
		tmpDir, cleanup, err := skills.CloneRepo(s)
		if err != nil {
			return "", nil, err
		}
		return tmpDir, cleanup, nil
	default:
		return "", nil, fmt.Errorf("--source requires a GitHub repo (owner/repo) or local path, got %s", source)
	}
}

// installAll installs every skill from the source directory to the target.
func installAll(installer *skills.Installer, sourceDir, parentDir string) error {
	skillsPath := filepath.Join(sourceDir, ".agents", "skills")
	if _, err := os.Stat(skillsPath); os.IsNotExist(err) {
		// Fall back to source root (mirrors InstallFromGitHub behavior)
		skillsPath = sourceDir
	}

	entries, err := os.ReadDir(skillsPath)
	if err != nil {
		return fmt.Errorf("reading skills directory %q: %w", skillsPath, err)
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if err := installer.Install(entry.Name(), parentDir); err != nil {
			fmt.Fprintf(os.Stderr, "warning: %v\n", err)
			continue
		}
		count++
	}

	if count == 0 {
		return fmt.Errorf("no skills found in %s", skillsPath)
	}

	fmt.Printf("Installed %d skills to %s\n", count, filepath.Join(parentDir, "skills"))
	return nil
}
