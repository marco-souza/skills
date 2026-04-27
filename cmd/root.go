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
)

// rootCmd represents the base command.
var rootCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage AI agent skills",
	Long: `Skills CLI - A tool for managing AI agent skills.

Use this CLI to list, create, validate, install, and search
skills in your projects or from remote repositories.`,
	Version:         fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		resolveConfig(cmd)
	},
}

// Execute runs the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("root", "r", "", "Root directory for skills operations")
	rootCmd.PersistentFlags().StringP("repo", "", "", "Remote GitHub repo (owner/repo or full URL)")
}

// resolveConfig loads persistent config and applies it as flag defaults
// when the user hasn't explicitly set them.
func resolveConfig(cmd *cobra.Command) {
	cfg, err := config.Load()
	if err != nil {
		return // silently use hardcoded defaults
	}
	if !cmd.Flags().Changed("root") {
		cmd.Flags().Set("root", cfg.DefaultRoot)
	}
	if !cmd.Flags().Changed("repo") {
		cmd.Flags().Set("repo", cfg.DefaultRepo)
	}
}
