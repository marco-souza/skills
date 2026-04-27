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

		root, _ := cmd.Flags().GetString("root")
		path := root
		if len(args) > 0 {
			path = args[0]
		}

		// Check if path points to a specific SKILL.md
		if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
			return validateSingle(path, fix)
		}

		validator := &skills.Validator{Fix: fix}
		results, err := validator.ValidateAll(path)
		if err != nil {
			return err
		}

		validCount := 0
		for _, r := range results {
			name := filepath.Base(filepath.Dir(r.SkillPath))
			if r.Valid {
				color.Green("✓ %s", name)
				validCount++
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

		fmt.Printf("\n%d/%d skills valid\n", validCount, len(results))
		if validCount < len(results) {
			os.Exit(1)
		}
		return nil
	},
}

func validateSingle(path string, fix bool) error {
	validator := &skills.Validator{Fix: fix}
	result := validator.Validate(path)

	name := filepath.Base(filepath.Dir(path))
	if result.Valid {
		color.Green("✓ %s", name)
	} else {
		color.Red("✗ %s", name)
		for _, e := range result.Errors {
			color.Red("  error:   %s", e)
		}
	}
	for _, w := range result.Warnings {
		color.Yellow("  warning: %s", w)
	}
	for _, f := range result.Fixed {
		color.Cyan("  fixed:   %s", f)
	}

	if !result.Valid {
		os.Exit(1)
	}
	return nil
}
