package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("format", "f", "table", "Output format (table, json)")
	listCmd.Flags().StringP("category", "c", "", "Filter by category")
}

var listCmd = &cobra.Command{
	Use:     "list [path]",
	Aliases: []string{"ls"},
	Short:   "List available skills",
	Long:  `List skills from the local directory or a remote GitHub repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		category, _ := cmd.Flags().GetString("category")

		root, _ := cmd.Flags().GetString("root")
		path := root
		if len(args) > 0 {
			path = args[0]
		}

		loader := skills.NewLoader(path)
		list, err := loader.LoadAll()
		if err != nil {
			return err
		}

		list = skills.FilterSkills(list, category, nil)

		if format == "json" {
			return json.NewEncoder(os.Stdout).Encode(list)
		}

		if len(list) == 0 {
			fmt.Println("No skills found.")
			return nil
		}

		printSkillsTable(list)
		fmt.Printf("\n%d skills found\n", len(list))
		return nil
	},
}

func printSkillsTable(all []*skills.Skill) {
	// Calculate column widths
	nameW := 4
	catW := 8
	descW := 11

	for _, s := range all {
		if len(s.Name) > nameW {
			nameW = len(s.Name)
		}
		if len(s.Category) > catW {
			catW = len(s.Category)
		}
		if len(s.Description) > descW {
			descW = len(s.Description)
		}
	}

	// Cap description width
	if descW > 60 {
		descW = 60
	}

	// Print header
	fmt.Printf("%-*s  %-*s  %-*s\n", nameW, "NAME", catW, "CATEGORY", descW, "DESCRIPTION")
	fmt.Printf("%-*s  %-*s  %-*s\n", nameW, "----", catW, "--------", descW, "-----------")

	for _, s := range all {
		desc := s.Description
		if len(desc) > descW {
			desc = desc[:descW-3] + "..."
		}
		cat := s.Category
		if cat == "" {
			cat = "-"
		}
		fmt.Printf("%-*s  %-*s  %s\n", nameW, s.Name, catW, cat, desc)
	}
}
