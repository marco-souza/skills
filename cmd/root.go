package cmd

import (
	"fmt"
	"os"

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
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
}

// Execute runs the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("root", "r", ".", "Root directory for skills operations")
	rootCmd.PersistentFlags().StringP("repo", "", "marco-souza/skills", "Remote GitHub repo (owner/repo or full URL)")
}
