package cmd

import (
	"fmt"
	"io"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
)

var validateCmd = newValidateCmd()

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:            "validate [path]",
		Short:          "Validate skills in a directory",
		SilenceErrors: true,
		SilenceUsage:  true,
		Long: `Validate all skills in the .agents/skills directory.

Checks each skill's SKILL.md file for valid frontmatter, required fields,
and format constraints. Reports pass/fail status for each skill.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			source, err := cmd.Flags().GetString("source")
			if err != nil {
				return fmt.Errorf("internal error reading --source flag: %w", err)
			}

			return validateSkills(cmd.OutOrStdout(), source)
		},
	}
	cmd.Flags().StringP("source", "s", "", "Source for skills: local path to project root (defaults to current directory)")
	return cmd
}

func validateSkills(w io.Writer, source string) error {
	sourceDir, cleanup, err := skills.ResolveSourceDir(source, cfg.DefaultSource)
	if err != nil {
		return err
	}
	if cleanup != nil {
		defer cleanup()
	}

	loader := skills.NewLoader(sourceDir)
	skillList, err := loader.LoadAll()
	if err != nil {
		return err
	}

	if len(skillList) == 0 {
		fmt.Fprintln(w, "no skills found")
		return nil
	}

	var invalidCount int

	for _, s := range skillList {
		if err := s.Validate(); err != nil {
			fmt.Fprintf(w, "FAIL  %s\n", s.Name)
			fmt.Fprintf(w, "      %v\n", err)
			invalidCount++
		} else {
			fmt.Fprintf(w, "PASS  %s\n", s.Name)
		}
	}

	fmt.Fprintf(w, "\n%d skill(s) validated, %d passed, %d failed\n",
		len(skillList), len(skillList)-invalidCount, invalidCount)

	if invalidCount > 0 {
		return &validationExitError{count: invalidCount}
	}

	return nil
}

type validationExitError struct {
	count int
}

func (e *validationExitError) Error() string {
	return fmt.Sprintf("%d skill(s) failed validation", e.count)
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
