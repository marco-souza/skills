package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringP("format", "f", "table", "Output format (table, json)")
	searchCmd.Flags().StringP("tag", "", "", "Filter by tag")
}

var searchCmd = &cobra.Command{
	Use:     "search <query>",
	Aliases: []string{"s"},
	Short:   "Search skills by name, description, or content",
	Long:  `Search for skills by querying name, description, or markdown content. Supports tag filtering.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		format, _ := cmd.Flags().GetString("format")
		tag, _ := cmd.Flags().GetString("tag")
		repo, _ := cmd.Flags().GetString("repo")
		root, _ := cmd.Flags().GetString("root")

		return skills.WithRemoteRepo(repo, root, func(root string) error {
			loader := skills.NewLoader(root)
			all, err := loader.LoadAll()
			if err != nil {
				return err
			}

			var results []*skills.Skill
			if tag != "" {
				results = skills.FilterSkills(all, "", []string{tag})
			} else {
				results = skills.SearchSkills(all, query)
			}

			if len(results) == 0 {
				fmt.Println("No skills found.")
				return nil
			}

			if format == "json" {
				return json.NewEncoder(os.Stdout).Encode(results)
			}

			printSearchTable(results)
			fmt.Printf("\n%d skills found\n", len(results))
			return nil
		})
	},
}

func printSearchTable(all []*skills.Skill) {
	nameW := 4
	catW := 8
	tagW := 4
	descW := 11

	for _, s := range all {
		if len(s.Name) > nameW {
			nameW = len(s.Name)
		}
		if len(s.Category) > catW {
			catW = len(s.Category)
		}
		tags := formatTags(s.Tags)
		if len(tags) > tagW {
			tagW = len(tags)
		}
		if len(s.Description) > descW {
			descW = len(s.Description)
		}
	}

	if descW > 50 {
		descW = 50
	}

	fmt.Printf("%-*s  %-*s  %-*s  %-*s\n", nameW, "NAME", catW, "CATEGORY", tagW, "TAGS", descW, "DESCRIPTION")
	fmt.Printf("%-*s  %-*s  %-*s  %-*s\n", nameW, "----", catW, "--------", tagW, "----", descW, "-----------")

	for _, s := range all {
		desc := s.Description
		if len(desc) > descW {
			desc = desc[:descW-3] + "..."
		}
		cat := s.Category
		if cat == "" {
			cat = "-"
		}
		tags := formatTags(s.Tags)
		if tags == "" {
			tags = "-"
		}
		fmt.Printf("%-*s  %-*s  %-*s  %s\n", nameW, s.Name, catW, cat, tagW, tags, desc)
	}
}

func formatTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	result := tags[0]
	for _, t := range tags[1:] {
		result += ", " + t
	}
	return result
}
