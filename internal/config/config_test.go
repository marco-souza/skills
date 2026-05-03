package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.DefaultSource != DefaultSource {
		t.Errorf("expected default_source %q, got %q", DefaultSource, cfg.DefaultSource)
	}
	if cfg.DefaultRoot != DefaultRoot {
		t.Errorf("expected default_root %q, got %q", DefaultRoot, cfg.DefaultRoot)
	}
}

func TestLoadSave(t *testing.T) {
	// Override config dir for test
	testDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", testDir)

	cfg := Default()
	cfg.DefaultSource = "test-org/test-repo"
	cfg.DefaultRoot = "/test/path"

	if err := cfg.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.DefaultSource != "test-org/test-repo" {
		t.Errorf("expected source %q, got %q", "test-org/test-repo", loaded.DefaultSource)
	}
	if loaded.DefaultRoot != "/test/path" {
		t.Errorf("expected root %q, got %q", "/test/path", loaded.DefaultRoot)
	}
}

func TestLoadDefaultsWhenMissing(t *testing.T) {
	// Use a temp dir with no config file
	testDir := t.TempDir()
	home := filepath.Join(testDir, "home")
	_ = os.MkdirAll(home, 0755)
	t.Setenv("HOME", home)

	// Force config to look in a non-existent directory
	// by using a fresh temp dir for XDG
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(testDir, "xdg"))

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if cfg.DefaultSource != DefaultSource {
		t.Errorf("expected default source %q, got %q", DefaultSource, cfg.DefaultSource)
	}
	if cfg.DefaultRoot != DefaultRoot {
		t.Errorf("expected default root %q, got %q", DefaultRoot, cfg.DefaultRoot)
	}
}
