package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a project with .agents structure",
	Long:  `Create the .agents/skills directory structure in the specified or current directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root, _ := cmd.Flags().GetString("root")
		path := root
		if len(args) > 0 {
			path = args[0]
		}

		// Resolve to absolute path for display
		abs, err := filepath.Abs(path)
		if err != nil {
			abs = path
		}

		skillsDir := filepath.Join(abs, ".agents", "skills")
		if err := os.MkdirAll(skillsDir, 0755); err != nil {
			return fmt.Errorf("creating skills directory: %w", err)
		}

		// Create AGENTS.md skeleton if it doesn't exist
		agentsFile := filepath.Join(abs, ".agents", "AGENTS.md")
		if _, err := os.Stat(agentsFile); os.IsNotExist(err) {
			content := `# AGENTS.md - AI Agent Skills

This directory contains AI agent skills for this project.

## Structure

\` + "``" + `
.agents/
└── skills/          # Skill definitions (SKILL.md files)
    └── my-skill/
        └── SKILL.md
\` + "``" + `

## Usage

Use the \` + "`" + `skills\` + "`" + ` CLI to manage skills:

\` + "``" + `bash
skills list          # List available skills
skills add my-skill  # Create a new skill
skills validate      # Validate all skills
\` + "``" + `
`
			if err := os.WriteFile(agentsFile, []byte(content), 0644); err != nil {
				return fmt.Errorf("creating AGENTS.md: %w", err)
			}
		}

		fmt.Printf("Initialized skills structure in %s\n", skillsDir)
		return nil
	},
}
