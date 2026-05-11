// Package cmd defines the Cobra CLI commands for the Skills tool.
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
	err     error
)

var rootCmd = &cobra.Command{
	Use:   "skills",
	Short: "Manage AI agent skills",
	Long: `Skills CLI — manage AI agent skills in your projects.

Use this CLI to list, install, uninstall, and create skills
in .agents/skills/ directories.`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	PersistentPreRunE: func(*cobra.Command, []string) error {
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

	// Register subcommands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(configCmd)

	// Register subcommand flags
	installCmd.Flags().StringP("target", "t", "", "Target project directory")
	installCmd.Flags().Bool("all", false, "Install all skills from the source")
	uninstallCmd.Flags().StringP("target", "t", "", "Target project directory")

	// Register config subcommands
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configListCmd)
}
