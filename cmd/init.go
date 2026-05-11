package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const agentsTemplate = `# AGENTS.md - AI Agent Skills

This directory contains AI agent skills for this project.

## Structure

%s
.agents/
└── skills/          # Skill definitions (SKILL.md files)
    └── my-skill/
        └── SKILL.md
%s

## Usage

Use the %sskills%s CLI to manage skills:

%s
skills list              # List available skills
skills add my-skill      # Create a new skill
skills install my-skill  # Install a skill
skills uninstall my-skill # Remove a skill
%s
`

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize a project with .agents structure",
	Long:  `Create the .agents/skills directory structure in the specified or current directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		abs, err := filepath.Abs(path)
		if err != nil {
			abs = path
		}

		skillsDir := filepath.Join(abs, ".agents", "skills")
		if err := os.MkdirAll(skillsDir, 0755); err != nil {
			return fmt.Errorf("creating skills directory: %w", err)
		}

		agentsFile := filepath.Join(abs, ".agents", "AGENTS.md")
		if _, err := os.Stat(agentsFile); os.IsNotExist(err) {
			fence := "```"
			code := "`"
			content := fmt.Sprintf(agentsTemplate, fence, fence, code, code, fence, fence)
			if err := os.WriteFile(agentsFile, []byte(content), 0644); err != nil {
				return fmt.Errorf("creating AGENTS.md: %w", err)
			}
		}

		fmt.Printf("Initialized skills structure in %s\n", skillsDir)
		return nil
	},
}
