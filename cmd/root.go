package cmd

import (
	"fmt"
	"os"

	"github.com/marco-souza/skills/internal/config"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage AI agent skills",
	Long: `Skills CLI — manage AI agent skills in your projects.

Use this CLI to list, install, uninstall, and create skills
in .agents/skills/ directories.`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	PersistentPreRunE: func(*cobra.Command, []string) error {
		var err error
		cfg, err = config.Load()
		return err
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("source", "s", "", "Source for skills: GitHub repo (owner/repo) or local path")
}
