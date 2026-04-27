package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
)

const defaultSkillTemplate = `---
name: {{.Name}}
description: >
  [What this skill does in one sentence].
  Use when: [positive triggers].
  Do NOT use when: [negative triggers].
---

# {{.Title}}

[Brief introduction]

## Prerequisites

[Any requirements before using the skill]

## When to Use

- Use case 1
- Use case 2

## When NOT to Use

- Anti-pattern 1
- Anti-pattern 2

## Instructions

[Core instructions here]

## Examples

[Concrete, copy-pasteable examples]

## Edge Cases / Troubleshooting

[Common issues and solutions]

## Best Practices

[Tips for optimal usage]
`

type skillTemplateData struct {
	Name  string
	Title string
}

func toTitle(s string) string {
	words := strings.Split(strings.ReplaceAll(s, "-", " "), " ")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("from", "f", "", "Use a custom template file")
	addCmd.Flags().StringP("target", "d", "", "Target project directory (default: current)")
}

var addCmd = &cobra.Command{
	Use:     "add <name>",
	Aliases: []string{"a"},
	Short:   "Create a new skill from template",
	Long:  `Create a new skill directory with a SKILL.md file from the default or custom template.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		target, _ := cmd.Flags().GetString("target")
		fromTemplate, _ := cmd.Flags().GetString("from")

		if target == "" {
			target, _ = cmd.Flags().GetString("root")
		}

		skillsDir := skills.ResolveToSkillsDir(target)
		skillPath := filepath.Join(skillsDir, name)
		skillFile := filepath.Join(skillPath, "SKILL.md")

		if _, err := os.Stat(skillFile); err == nil {
			return fmt.Errorf("skill %q already exists at %s", name, skillFile)
		}

		if err := os.MkdirAll(skillPath, 0755); err != nil {
			return fmt.Errorf("creating skill directory: %w", err)
		}

		var tmpl *template.Template
		if fromTemplate != "" {
			var err error
			tmpl, err = template.ParseFiles(fromTemplate)
			if err != nil {
				return fmt.Errorf("parsing template: %w", err)
			}
		} else {
			tmpl = template.Must(template.New("skill").Parse(defaultSkillTemplate))
		}

		f, err := os.Create(skillFile)
		if err != nil {
			return fmt.Errorf("creating SKILL.md: %w", err)
		}
		defer f.Close()

		title := toTitle(name)
		data := skillTemplateData{Name: name, Title: title}

		if err := tmpl.Execute(f, data); err != nil {
			return fmt.Errorf("rendering template: %w", err)
		}

		fmt.Printf("Created skill %q at %s\n", name, skillFile)
		return nil
	},
}
