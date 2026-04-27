package cmd

import (
	"fmt"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringP("tag", "", "", "Filter by tag")
}

var searchCmd = &cobra.Command{
	Use:     "search <query>",
	Aliases: []string{"s"},
	Short:   "Search skills by name, description, or content",
	Long:    `Search for skills by querying name, description, or markdown content. Supports tag filtering.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
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

			out, err := yaml.Marshal(results)
			if err != nil {
				return err
			}
			fmt.Print(string(out))
			return nil
		})
	},
}
