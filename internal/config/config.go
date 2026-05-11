// Package config handles persistent CLI configuration stored at
// ~/.config/skills/config.yaml.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Default configuration values.
const (
	// DefaultSource is the default GitHub repository used as the skill source.
	DefaultSource = "marco-souza/skills"
	// DefaultRoot is the default project root directory.
	DefaultRoot = "."
)

// Config holds persistent CLI settings loaded from ~/.config/skills/config.yaml.
type Config struct {
	// DefaultSource is the default GitHub skill source (owner/repo format).
	DefaultSource string `yaml:"default_source"`
	// DefaultRoot is the default project root directory.
	DefaultRoot string `yaml:"default_root"`
}

// Default returns a new Config populated with factory default values.
func Default() *Config {
	return &Config{
		DefaultSource: DefaultSource,
		DefaultRoot:   DefaultRoot,
	}
}

// Dir returns the configuration directory path (~/.config/skills).
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}
	return filepath.Join(home, ".config", "skills"), nil
}

// Path returns the full path to the configuration file (~/.config/skills/config.yaml).
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// Load reads the configuration from disk, falling back to factory defaults if the file does not exist.
func Load() (*Config, error) {
	cfg := Default()

	path, err := Path()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: cannot determine config path: %v\n", err)
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// Apply defaults for empty fields
	if cfg.DefaultSource == "" {
		cfg.DefaultSource = DefaultSource
	}
	if cfg.DefaultRoot == "" {
		cfg.DefaultRoot = DefaultRoot
	}

	return cfg, nil
}

// Save writes the configuration to disk, creating the parent directory if it does not exist.
func (c *Config) Save() error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	path, err := Path()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
