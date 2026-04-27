// Package config handles persistent CLI configuration stored at
// ~/.config/skills/config.yaml.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Default values
const (
	DefaultRepo = "marco-souza/skills"
	DefaultRoot = "."
)

// Config holds persistent CLI settings.
type Config struct {
	DefaultRepo string `yaml:"default_repo"`
	DefaultRoot string `yaml:"default_root"`
}

// Default returns a Config with factory defaults.
func Default() *Config {
	return &Config{
		DefaultRepo: DefaultRepo,
		DefaultRoot: DefaultRoot,
	}
}

// Dir returns the config directory path (~/.config/skills).
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}
	return filepath.Join(home, ".config", "skills"), nil
}

// Path returns the full config file path.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// Load reads config from disk, falling back to defaults if the file doesn't exist.
func Load() (*Config, error) {
	cfg := Default()

	path, err := Path()
	if err != nil {
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
	if cfg.DefaultRepo == "" {
		cfg.DefaultRepo = DefaultRepo
	}
	if cfg.DefaultRoot == "" {
		cfg.DefaultRoot = DefaultRoot
	}

	return cfg, nil
}

// Save writes the config to disk, creating the directory if needed.
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
