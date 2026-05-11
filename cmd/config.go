package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage persistent CLI configuration",
	Long: `Manage persistent CLI configuration stored at ~/.config/skills/config.yaml.

Subcommands:
  skills config get <key>   Get a config value
  skills config set <k> <v> Set a config value
  skills config list        Show all config values`,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all config values",
	RunE: func(*cobra.Command, []string) error {
		out, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("marshaling config: %w", err)
		}
		fmt.Print(string(out))
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		switch args[0] {
		case "default_source":
			fmt.Println(cfg.DefaultSource)
		case "default_root":
			fmt.Println(cfg.DefaultRoot)
		default:
			return fmt.Errorf("unknown config key: %s (valid: default_source, default_root)", args[0])
		}
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Args:  cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		switch args[0] {
		case "default_source":
			cfg.DefaultSource = args[1]
		case "default_root":
			cfg.DefaultRoot = args[1]
		default:
			return fmt.Errorf("unknown config key: %s (valid: default_source, default_root)", args[0])
		}
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
		fmt.Printf("Set %s = %s\n", args[0], args[1])
		return nil
	},
}
