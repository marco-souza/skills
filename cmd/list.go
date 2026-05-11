package cmd

import (
	"fmt"
	"io"

	"github.com/marco-souza/skills/internal/skills"
	"github.com/spf13/cobra"
)

var listCmd = newListCmd()

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [path]",
		Aliases: []string{"ls"},
		Short:   "List available skills",
		Long:    `List skills from the local .agents/skills directory, a local path, or a remote GitHub repo.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			source, err := cmd.Flags().GetString("source")
			if err != nil {
				return fmt.Errorf("internal error reading --source flag: %w", err)
			}

			sourceDir, cleanup, err := skills.ResolveSourceDir(source, cfg.DefaultSource, defaultExecFunc)
			if err != nil {
				return err
			}
			if cleanup != nil {
				defer cleanup()
			}

			return listSkills(cmd.OutOrStdout(), sourceDir)
		},
	}
	cmd.Flags().StringP("source", "s", "", "Source for skills: GitHub repo (owner/repo) or local path")
	return cmd
}

func listSkills(w io.Writer, root string) error {
	loader := skills.NewLoader(root)
	sk, err := loader.LoadAll()
	if err != nil {
		return err
	}

	if len(sk) == 0 {
		fmt.Fprintln(w, "no skills found")
		return nil
	}

	for _, s := range sk {
		fmt.Fprintln(w, s.Name)
	}
	return nil
}


