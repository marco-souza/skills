package cmd

import (
	"fmt"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("category", "c", "", "Filter by category")
}

var listCmd = &cobra.Command{
	Use:     "list [path]",
	Aliases: []string{"ls"},
	Short:   "List available skills",
	Long:    `List skills from the local directory or a remote GitHub repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		category, _ := cmd.Flags().GetString("category")
		repo, _ := cmd.Flags().GetString("repo")
		root, _ := cmd.Flags().GetString("root")
		path := root
		if len(args) > 0 {
			path = args[0]
		}

		return skills.WithRemoteRepo(repo, path, func(root string) error {
			loader := skills.NewLoader(root)
			list, err := loader.LoadAll()
			if err != nil {
				return err
			}

			list = skills.FilterSkills(list, category, nil)

			out, err := yaml.Marshal(list)
			if err != nil {
				return err
			}
			fmt.Print(string(out))
			return nil
		})
	},
}
