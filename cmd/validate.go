package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().Bool("fix", false, "Auto-fix issues where possible")
}

var validateCmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate skill format and structure",
	Long:  `Validate all skills or a specific skill for correct format, required fields, and naming conventions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fix, _ := cmd.Flags().GetBool("fix")
		repo, _ := cmd.Flags().GetString("repo")
		root, _ := cmd.Flags().GetString("root")
		path := root
		if len(args) > 0 {
			path = args[0]
		}

		return skills.WithRemoteRepo(repo, path, func(root string) error {
			// Check if path points to a specific SKILL.md
			if fi, err := os.Stat(root); err == nil && !fi.IsDir() {
				return validateSingle(root, fix)
			}

			validator := &skills.Validator{Fix: fix}
			results, err := validator.ValidateAll(root)
			if err != nil {
				return err
			}

			return printValidateResults(results)
		})
	},
}

func validateSingle(path string, fix bool) error {
	validator := &skills.Validator{Fix: fix}
	result := validator.Validate(path)
	return printValidateResults([]skills.Result{result})
}

func printResult(r skills.Result) {
	name := filepath.Base(filepath.Dir(r.SkillPath))
	if r.Valid {
		color.Green("✓ %s", name)
	} else {
		color.Red("✗ %s", name)
		for _, e := range r.Errors {
			color.Red("  error:   %s", e)
		}
	}
	for _, w := range r.Warnings {
		color.Yellow("  warning: %s", w)
	}
	for _, f := range r.Fixed {
		color.Cyan("  fixed:   %s", f)
	}
}

func printValidateResults(results []skills.Result) error {
	validCount := 0
	for _, r := range results {
		printResult(r)
		if r.Valid {
			validCount++
		}
	}

	fmt.Printf("\n%d/%d skills valid\n", validCount, len(results))
	if validCount < len(results) {
		os.Exit(1)
	}
	return nil
}
