package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/marco-souza/skills/internal/config"
	"github.com/marco-souza/skills/internal/skills"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func setupSourceSkill(t *testing.T, src, dirName, name, description string) {
	t.Helper()
	skillDir := filepath.Join(src, skills.DefaultSkillsDir, skills.SkillsSubDir, dirName)
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("creating skill dir: %v", err)
	}
	content := fmt.Sprintf("---\nname: %s\ndescription: %s\n---\n\n# %s\n", name, description, name)
	os.WriteFile(filepath.Join(skillDir, skills.SkillFileName), []byte(content), 0644)
}

func setupCfg(t *testing.T) {
	t.Helper()
	cfg = config.Default()
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %v failed: %v", args, err)
	}
}

func createBareRepoWithSkills(t *testing.T) string {
	t.Helper()
	tempDir := t.TempDir()
	workDir := filepath.Join(tempDir, "work")
	bareRepo := filepath.Join(tempDir, "bare.git")

	os.MkdirAll(workDir, 0755)
	runGit(t, workDir, "init")
	runGit(t, workDir, "config", "user.email", "test@test.com")
	runGit(t, workDir, "config", "user.name", "Test")

	// Add two skills
	setupSourceSkill(t, workDir, "skill-alpha", "skill-alpha", "Alpha skill for testing.")
	setupSourceSkill(t, workDir, "skill-beta", "skill-beta", "Beta skill for testing.")

	runGit(t, workDir, "add", ".")
	runGit(t, workDir, "commit", "-m", "init")
	runGit(t, tempDir, "clone", "--bare", workDir, bareRepo)

	return bareRepo
}

// ── toTitleCase (pure function) ──────────────────────────────────────────────

func TestToTitleCase(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"single word", "hello", "Hello"},
		{"single char", "a", "A"},
		{"multi word hyphenated", "git-commit-formatter", "Git Commit Formatter"},
		{"already titled", "My Skill", "My Skill"},
		{"hyphenated single", "my-skill", "My Skill"},
		{"multiple hyphens", "a-b-c-d", "A B C D"},
		{"mixed case input", "GIT-COMMIT", "GIT COMMIT"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := toTitleCase(tc.input)
			if got != tc.want {
				t.Errorf("toTitleCase(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// ── listSkills ───────────────────────────────────────────────────────────────

func TestListSkills(t *testing.T) {
	t.Parallel()

	t.Run("no skills directory", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		err := listSkills(dir)
		if err == nil {
			t.Fatal("expected error for missing skills directory")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected 'not found' in error, got: %v", err)
		}
	})

	t.Run("empty skills directory", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		skillsDir := filepath.Join(dir, skills.DefaultSkillsDir, skills.SkillsSubDir)
		os.MkdirAll(skillsDir, 0755)

		var buf bytes.Buffer
		old := os.Stdout
		os.Stdout = &writeSyncer{Writer: &buf}
		defer func() { os.Stdout = old }()

		err := listSkills(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if buf.String() != "no skills found\n" {
			t.Errorf("expected 'no skills found', got %q", buf.String())
		}
	})

	t.Run("lists skill names", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		setupSourceSkill(t, dir, "my-skill", "my-skill", "A test skill.")
		setupSourceSkill(t, dir, "other-skill", "other-skill", "Another skill.")

		var buf bytes.Buffer
		old := os.Stdout
		os.Stdout = &writeSyncer{Writer: &buf}
		defer func() { os.Stdout = old }()

		err := listSkills(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		if len(lines) != 2 {
			t.Fatalf("expected 2 lines, got %d: %q", len(lines), buf.String())
		}
		// Order depends on filesystem; just check both are present
		found := map[string]bool{}
		for _, l := range lines {
			found[strings.TrimSpace(l)] = true
		}
		if !found["my-skill"] || !found["other-skill"] {
			t.Errorf("expected both skill names, got: %q", buf.String())
		}
	})

	t.Run("skips non-directories", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		skillsDir := filepath.Join(dir, skills.DefaultSkillsDir, skills.SkillsSubDir)
		os.MkdirAll(skillsDir, 0755)
		// Create a file (not a directory) in skills dir
		os.WriteFile(filepath.Join(skillsDir, "NOT-A-DIR"), []byte("junk"), 0644)
		// Create a real skill
		setupSourceSkill(t, dir, "real-skill", "real-skill", "Real skill.")

		var buf bytes.Buffer
		old := os.Stdout
		os.Stdout = &writeSyncer{Writer: &buf}
		defer func() { os.Stdout = old }()

		err := listSkills(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(buf.String(), "real-skill") {
			t.Errorf("expected 'real-skill' in output, got %q", buf.String())
		}
	})
}

// writeSyncer wraps an io.Writer to satisfy *os.File interface for stdout swap.
type writeSyncer struct {
	*bytes.Buffer
}

func (w *writeSyncer) Sync() error { return nil }
func (w *writeSyncer) Fd() uintptr { return 0 }

// ── init command ─────────────────────────────────────────────────────────────

func TestInitCommand(t *testing.T) {
	t.Parallel()

	t.Run("creates .agents/skills directory", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()

		c := initCmd
		c.SetArgs([]string{dir})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("init failed: %v", err)
		}

		skillsDir := filepath.Join(dir, ".agents", "skills")
		if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
			t.Fatal("skills directory was not created")
		}
	})

	t.Run("creates AGENTS.md", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()

		c := initCmd
		c.SetArgs([]string{dir})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("init failed: %v", err)
		}

		agentsFile := filepath.Join(dir, ".agents", "AGENTS.md")
		if _, err := os.Stat(agentsFile); os.IsNotExist(err) {
			t.Fatal("AGENTS.md was not created")
		}

		data, _ := os.ReadFile(agentsFile)
		content := string(data)
		if !strings.Contains(content, "AGENTS.md") {
			t.Error("AGENTS.md has unexpected content")
		}
		if !strings.Contains(content, "skills list") {
			t.Error("AGENTS.md missing usage examples")
		}
	})

	t.Run("does not overwrite existing AGENTS.md", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()

		agentsFile := filepath.Join(dir, ".agents", "AGENTS.md")
		os.MkdirAll(filepath.Dir(agentsFile), 0755)
		original := "# My custom AGENTS.md\n"
		os.WriteFile(agentsFile, []byte(original), 0644)

		c := initCmd
		c.SetArgs([]string{dir})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("init failed: %v", err)
		}

		data, _ := os.ReadFile(agentsFile)
		if string(data) != original {
			t.Error("existing AGENTS.md was overwritten")
		}
	})

	t.Run("defaults to current directory", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		cwd, _ := os.Getwd()
		_ = os.Chdir(dir)
		defer os.Chdir(cwd)

		c := initCmd
		c.SetArgs([]string{})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("init failed: %v", err)
		}

		skillsDir := filepath.Join(dir, ".agents", "skills")
		if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
			t.Fatal("skills directory was not created in cwd")
		}
	})
}

// ── config commands ──────────────────────────────────────────────────────────

func TestConfigGetCmd(t *testing.T) {
	t.Parallel()

	t.Run("get default_source", func(t *testing.T) {
		t.Parallel()
		cfg = &config.Config{DefaultSource: "my-org/my-repo", DefaultRoot: "/root"}

		var buf bytes.Buffer
		c := configGetCmd
		c.SetArgs([]string{"default_source"})
		c.SetOut(&buf)
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("config get failed: %v", err)
		}
		got := strings.TrimSpace(buf.String())
		if got != "my-org/my-repo" {
			t.Errorf("expected 'my-org/my-repo', got %q", got)
		}
	})

	t.Run("get default_root", func(t *testing.T) {
		t.Parallel()
		cfg = &config.Config{DefaultSource: "my-org/my-repo", DefaultRoot: "/custom/root"}

		var buf bytes.Buffer
		c := configGetCmd
		c.SetArgs([]string{"default_root"})
		c.SetOut(&buf)
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("config get failed: %v", err)
		}
		got := strings.TrimSpace(buf.String())
		if got != "/custom/root" {
			t.Errorf("expected '/custom/root', got %q", got)
		}
	})

	t.Run("unknown key", func(t *testing.T) {
		t.Parallel()
		cfg = config.Default()

		c := configGetCmd
		c.SetArgs([]string{"unknown_key"})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err == nil {
			t.Fatal("expected error for unknown key")
		}
		if !strings.Contains(err.Error(), "unknown config key") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing argument", func(t *testing.T) {
		t.Parallel()
		cfg = config.Default()

		c := configGetCmd
		c.SetArgs([]string{})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err == nil {
			t.Fatal("expected error for missing argument")
		}
	})
}

func TestConfigSetCmd(t *testing.T) {
	t.Parallel()

	t.Run("set default_source", func(t *testing.T) {
		t.Parallel()
		// Use a temp HOME so config writes to a temp location
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		cfg = &config.Config{DefaultSource: "old/repo", DefaultRoot: "."}

		var buf bytes.Buffer
		c := configSetCmd
		c.SetArgs([]string{"default_source", "new-org/new-repo"})
		c.SetOut(&buf)
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("config set failed: %v", err)
		}

		if !strings.Contains(buf.String(), "Set default_source = new-org/new-repo") {
			t.Errorf("unexpected output: %q", buf.String())
		}
		if cfg.DefaultSource != "new-org/new-repo" {
			t.Errorf("cfg.DefaultSource = %q, want 'new-org/new-repo'", cfg.DefaultSource)
		}
	})

	t.Run("set default_root", func(t *testing.T) {
		t.Parallel()
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		cfg = &config.Config{DefaultSource: "my/repo", DefaultRoot: "."}

		c := configSetCmd
		c.SetArgs([]string{"default_root", "/new/root"})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("config set failed: %v", err)
		}
		if cfg.DefaultRoot != "/new/root" {
			t.Errorf("cfg.DefaultRoot = %q, want '/new/root'", cfg.DefaultRoot)
		}
	})

	t.Run("unknown key", func(t *testing.T) {
		t.Parallel()
		cfg = config.Default()

		c := configSetCmd
		c.SetArgs([]string{"bogus_key", "value"})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err == nil {
			t.Fatal("expected error for unknown key")
		}
		if !strings.Contains(err.Error(), "unknown config key") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestConfigListCmd(t *testing.T) {
	t.Parallel()

	t.Run("lists all config values as YAML", func(t *testing.T) {
		t.Parallel()
		cfg = &config.Config{DefaultSource: "org/repo", DefaultRoot: "/root"}

		var buf bytes.Buffer
		c := configListCmd
		c.SetArgs([]string{})
		c.SetOut(&buf)
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("config list failed: %v", err)
		}
		out := buf.String()
		if !strings.Contains(out, "default_source") {
			t.Error("output missing default_source")
		}
		if !strings.Contains(out, "default_root") {
			t.Error("output missing default_root")
		}
	})
}

// ── install command ──────────────────────────────────────────────────────────

func TestInstallCommand(t *testing.T) {
	t.Parallel()

	t.Run("install single skill", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		// Create source with skills
		src := t.TempDir()
		setupSourceSkill(t, src, "my-skill", "my-skill", "A test skill.")

		target := t.TempDir()

		c := installCmd
		c.SetArgs([]string{"my-skill", "-s", src, "-t", target})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("install failed: %v", err)
		}

		dest := filepath.Join(target, ".agents", "skills", "my-skill", skills.SkillFileName)
		if _, err := os.Stat(dest); os.IsNotExist(err) {
			t.Fatal("skill was not installed")
		}
	})

	t.Run("install multiple skills", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		src := t.TempDir()
		setupSourceSkill(t, src, "skill-a", "skill-a", "Skill A.")
		setupSourceSkill(t, src, "skill-b", "skill-b", "Skill B.")

		target := t.TempDir()

		c := installCmd
		c.SetArgs([]string{"skill-a", "skill-b", "-s", src, "-t", target})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("install failed: %v", err)
		}

		for _, name := range []string{"skill-a", "skill-b"} {
			dest := filepath.Join(target, ".agents", "skills", name, skills.SkillFileName)
			if _, err := os.Stat(dest); os.IsNotExist(err) {
				t.Fatalf("skill %q was not installed", name)
			}
		}
	})

	t.Run("install all skills", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		src := t.TempDir()
		setupSourceSkill(t, src, "alpha", "alpha", "Alpha skill.")
		setupSourceSkill(t, src, "beta", "beta", "Beta skill.")
		setupSourceSkill(t, src, "gamma", "gamma", "Gamma skill.")

		target := t.TempDir()

		c := installCmd
		c.SetArgs([]string{"--all", "-s", src, "-t", target})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("install --all failed: %v", err)
		}

		for _, name := range []string{"alpha", "beta", "gamma"} {
			dest := filepath.Join(target, ".agents", "skills", name, skills.SkillFileName)
			if _, err := os.Stat(dest); os.IsNotExist(err) {
				t.Fatalf("skill %q was not installed with --all", name)
			}
		}
	})

	t.Run("install all with scripts", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		src := t.TempDir()
		skillDir := setupSourceSkillDir(t, src, "scripted-skill")
		os.WriteFile(filepath.Join(skillDir, skills.SkillFileName), []byte(`---
name: scripted-skill
description: Has a script.
metadata:
  scripts:
    - scripts/helper.sh
---

# Scripted Skill
`), 0644)

		scriptDir := filepath.Join(skillDir, "scripts")
		os.MkdirAll(scriptDir, 0755)
		os.WriteFile(filepath.Join(scriptDir, "helper.sh"), []byte("#!/bin/bash\n"), 0755)

		target := t.TempDir()

		c := installCmd
		c.SetArgs([]string{"scripted-skill", "-s", src, "-t", target})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("install failed: %v", err)
		}

		// Verify skill
		destSkill := filepath.Join(target, ".agents", "skills", "scripted-skill", skills.SkillFileName)
		if _, err := os.Stat(destSkill); os.IsNotExist(err) {
			t.Fatal("skill was not installed")
		}
		// Verify script
		destScript := filepath.Join(target, ".agents", "scripts", "helper.sh")
		if _, err := os.Stat(destScript); os.IsNotExist(err) {
			t.Fatal("script was not installed")
		}
	})

	t.Run("install all skips non-directories", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		src := t.TempDir()
		skillsDir := filepath.Join(src, skills.DefaultSkillsDir, skills.SkillsSubDir)
		os.MkdirAll(skillsDir, 0755)
		// Create a file that's not a directory
		os.WriteFile(filepath.Join(skillsDir, "README.md"), []byte("not a skill"), 0644)
		// Create a real skill
		setupSourceSkill(t, src, "real-skill", "real-skill", "A real skill.")

		target := t.TempDir()

		c := installCmd
		c.SetArgs([]string{"--all", "-s", src, "-t", target})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("install --all failed: %v", err)
		}

		// Only the real skill should be installed
		dest := filepath.Join(target, ".agents", "skills", "real-skill", skills.SkillFileName)
		if _, err := os.Stat(dest); os.IsNotExist(err) {
			t.Fatal("real skill was not installed")
		}
	})

	t.Run("rejects skill names with --all", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		c := installCmd
		c.SetArgs([]string{"some-skill", "--all"})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err == nil {
			t.Fatal("expected error when specifying skill names with --all")
		}
		if !strings.Contains(err.Error(), "cannot specify skill names with --all") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("rejects missing skill names without --all", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		c := installCmd
		c.SetArgs([]string{})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err == nil {
			t.Fatal("expected error when no skill names and no --all")
		}
		if !strings.Contains(err.Error(), "requires at least one skill name") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("install nonexistent skill", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		src := t.TempDir()
		skillsDir := filepath.Join(src, skills.DefaultSkillsDir, skills.SkillsSubDir)
		os.MkdirAll(skillsDir, 0755)

		target := t.TempDir()

		c := installCmd
		c.SetArgs([]string{"nonexistent", "-s", src, "-t", target})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err == nil {
			t.Fatal("expected error for nonexistent skill")
		}
	})

	t.Run("install with --all from empty source", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		src := t.TempDir()
		// No .agents/skills directory at all
		skillsDir := filepath.Join(src, skills.DefaultSkillsDir, skills.SkillsSubDir)
		os.MkdirAll(skillsDir, 0755)

		target := t.TempDir()

		c := installCmd
		c.SetArgs([]string{"--all", "-s", src, "-t", target})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err == nil {
			t.Fatal("expected error for empty source with --all")
		}
		if !strings.Contains(err.Error(), "no skills found") {
			t.Errorf("expected 'no skills found' in error, got: %v", err)
		}
	})
}

// setupSourceSkillDir creates just the directory for a source skill (no SKILL.md).
func setupSourceSkillDir(t *testing.T, src, dirName string) string {
	t.Helper()
	skillDir := filepath.Join(src, skills.DefaultSkillsDir, skills.SkillsSubDir, dirName)
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("creating skill dir: %v", err)
	}
	return skillDir
}

// ── install via GitHub (local bare repo) ─────────────────────────────────────

func TestInstallViaGitHub(t *testing.T) {
	t.Parallel()

	t.Run("install from local bare repo", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		bareRepo := createBareRepoWithSkills(t)

		// Override execCommand to use local bare repo instead of SSH
		origExec := skills.ExecCommand
		defer func() { skills.ExecCommand = origExec }()

		skills.ExecCommand = func(name string, args ...string) *exec.Cmd {
			// Replace the SSH URL argument with the local bare repo path
			newArgs := make([]string, len(args))
			copy(newArgs, args)
			for i, arg := range newArgs {
				// The SSH URL is the last non-directory argument (position varies)
				if strings.HasSuffix(arg, ".git") || strings.Contains(arg, "bare") || arg == bareRepo {
					newArgs[i] = bareRepo
				}
			}
			// If we couldn't find it by matching, replace the clone source (5th arg for git clone --depth 1 <url> <dir>)
			// args format: clone --depth 1 <SSHURL> <repoDir>
			if len(newArgs) >= 5 {
				newArgs[3] = bareRepo
			}
			return exec.Command(name, newArgs...)
		}

		target := t.TempDir()

		c := installCmd
		c.SetArgs([]string{"--all", "-s", "test/test-repo", "-t", target})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err != nil {
			t.Fatalf("install from GitHub failed: %v", err)
		}

		// Verify both skills were installed
		for _, name := range []string{"skill-alpha", "skill-beta"} {
			dest := filepath.Join(target, ".agents", "skills", name, skills.SkillFileName)
			if _, err := os.Stat(dest); os.IsNotExist(err) {
				t.Fatalf("skill %q was not installed from GitHub source", name)
			}
		}
	})

	t.Run("GitHub clone failure returns error", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		origExec := skills.ExecCommand
		defer func() { skills.ExecCommand = origExec }()

		skills.ExecCommand = func(name string, args ...string) *exec.Cmd {
			return exec.Command("false")
		}

		target := t.TempDir()

		c := installCmd
		c.SetArgs([]string{"some-skill", "-s", "fail-test/fail-repo", "-t", target})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err == nil {
			t.Fatal("expected error from failed GitHub clone")
		}
	})
}

// ── list command ─────────────────────────────────────────────────────────────

func TestListCommand(t *testing.T) {
	t.Parallel()

	t.Run("list from local source", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		src := t.TempDir()
		setupSourceSkill(t, src, "skill-one", "skill-one", "First skill.")
		setupSourceSkill(t, src, "skill-two", "skill-two", "Second skill.")

		var buf bytes.Buffer
		c := listCmd
		c.SetArgs([]string{"-s", src})
		c.SetOut(&buf)
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("list failed: %v", err)
		}

		out := buf.String()
		if !strings.Contains(out, "skill-one") {
			t.Error("output missing skill-one")
		}
		if !strings.Contains(out, "skill-two") {
			t.Error("output missing skill-two")
		}
	})

	t.Run("list from empty source", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		src := t.TempDir()
		skillsDir := filepath.Join(src, skills.DefaultSkillsDir, skills.SkillsSubDir)
		os.MkdirAll(skillsDir, 0755)

		var buf bytes.Buffer
		c := listCmd
		c.SetArgs([]string{"-s", src})
		c.SetOut(&buf)
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("list failed: %v", err)
		}
		if !strings.Contains(buf.String(), "no skills found") {
			t.Errorf("expected 'no skills found', got %q", buf.String())
		}
	})

	t.Run("list from nonexistent source", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		c := listCmd
		c.SetArgs([]string{"-s", "/nonexistent/path/that/does/not/exist"})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err == nil {
			t.Fatal("expected error for nonexistent source")
		}
	})
}

// ── uninstall command ────────────────────────────────────────────────────────

func TestUninstallCommand(t *testing.T) {
	t.Parallel()

	t.Run("uninstall existing skill", func(t *testing.T) {
		t.Parallel()
		target := t.TempDir()
		skillsDir := filepath.Join(target, ".agents", "skills", "old-skill")
		os.MkdirAll(skillsDir, 0755)
		os.WriteFile(filepath.Join(skillsDir, skills.SkillFileName), []byte("---\nname: old-skill\ndescription: test\n---\n"), 0644)

		c := uninstallCmd
		c.SetArgs([]string{"old-skill", "-t", target})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("uninstall failed: %v", err)
		}

		if _, err := os.Stat(skillsDir); !os.IsNotExist(err) {
			t.Fatal("skill directory was not removed")
		}
	})

	t.Run("uninstall nonexistent skill", func(t *testing.T) {
		t.Parallel()
		target := t.TempDir()
		skillsDir := filepath.Join(target, ".agents", "skills")
		os.MkdirAll(skillsDir, 0755)

		c := uninstallCmd
		c.SetArgs([]string{"nonexistent", "-t", target})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err == nil {
			t.Fatal("expected error for nonexistent skill")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("uninstall multiple skills", func(t *testing.T) {
		t.Parallel()
		target := t.TempDir()
		for _, name := range []string{"skill-a", "skill-b"} {
			skillsDir := filepath.Join(target, ".agents", "skills", name)
			os.MkdirAll(skillsDir, 0755)
			os.WriteFile(filepath.Join(skillsDir, skills.SkillFileName), []byte("---\nname: "+name+"\ndescription: test\n---\n"), 0644)
		}

		c := uninstallCmd
		c.SetArgs([]string{"skill-a", "skill-b", "-t", target})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("uninstall failed: %v", err)
		}

		for _, name := range []string{"skill-a", "skill-b"} {
			dest := filepath.Join(target, ".agents", "skills", name)
			if _, err := os.Stat(dest); !os.IsNotExist(err) {
				t.Fatalf("skill %q was not removed", name)
			}
		}
	})

	t.Run("uninstall without skill names", func(t *testing.T) {
		t.Parallel()
		c := uninstallCmd
		c.SetArgs([]string{})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err == nil {
			t.Fatal("expected error for missing skill names")
		}
	})

	t.Run("uninstall with custom target dir", func(t *testing.T) {
		t.Parallel()
		target := t.TempDir()
		skillsDir := filepath.Join(target, "custom", "skills", "my-skill")
		os.MkdirAll(skillsDir, 0755)
		os.WriteFile(filepath.Join(skillsDir, skills.SkillFileName), []byte("---\nname: my-skill\ndescription: test\n---\n"), 0644)

		c := uninstallCmd
		c.SetArgs([]string{"my-skill", "-t", filepath.Join(target, "custom")})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("uninstall failed: %v", err)
		}

		if _, err := os.Stat(skillsDir); !os.IsNotExist(err) {
			t.Fatal("skill was not removed from custom target")
		}
	})
}

// ── add command ──────────────────────────────────────────────────────────────

func TestAddCommand(t *testing.T) {
	t.Parallel()

	t.Run("creates new skill", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		skillsDir := filepath.Join(dir, skills.DefaultSkillsDir, skills.SkillsSubDir)
		os.MkdirAll(skillsDir, 0755)

		// Change to temp dir so addCmd finds the skills dir
		cwd, _ := os.Getwd()
		_ = os.Chdir(dir)
		defer os.Chdir(cwd)

		c := addCmd
		c.SetArgs([]string{"my-new-skill"})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("add failed: %v", err)
		}

		skillFile := filepath.Join(skillsDir, "my-new-skill", skills.SkillFileName)
		if _, err := os.Stat(skillFile); os.IsNotExist(err) {
			t.Fatal("SKILL.md was not created")
		}

		data, _ := os.ReadFile(skillFile)
		content := string(data)
		if !strings.Contains(content, "name: my-new-skill") {
			t.Error("SKILL.md missing skill name")
		}
		if !strings.Contains(content, "# My New Skill") {
			t.Error("SKILL.md missing title")
		}
	})

	t.Run("rejects duplicate skill", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		skillsDir := filepath.Join(dir, skills.DefaultSkillsDir, skills.SkillsSubDir)
		skillDir := filepath.Join(skillsDir, "existing-skill")
		os.MkdirAll(skillDir, 0755)
		os.WriteFile(filepath.Join(skillDir, skills.SkillFileName), []byte("---\nname: existing-skill\ndescription: exists\n---\n"), 0644)

		cwd, _ := os.Getwd()
		_ = os.Chdir(dir)
		defer os.Chdir(cwd)

		c := addCmd
		c.SetArgs([]string{"existing-skill"})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err == nil {
			t.Fatal("expected error for duplicate skill")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing argument", func(t *testing.T) {
		t.Parallel()
		c := addCmd
		c.SetArgs([]string{})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err == nil {
			t.Fatal("expected error for missing argument")
		}
	})
}

// ── resolve source dir via commands ──────────────────────────────────────────

func TestResolveSourceDirViaCommands(t *testing.T) {
	t.Parallel()

	t.Run("local .agents found when source empty", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		dir := t.TempDir()
		setupSourceSkill(t, dir, "local-skill", "local-skill", "Local skill.")

		cwd, _ := os.Getwd()
		_ = os.Chdir(dir)
		defer os.Chdir(cwd)

		var buf bytes.Buffer
		c := listCmd
		c.SetArgs([]string{})
		c.SetOut(&buf)
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("list with local .agents failed: %v", err)
		}
		if !strings.Contains(buf.String(), "local-skill") {
			t.Errorf("expected 'local-skill' in output, got %q", buf.String())
		}
	})

	t.Run("explicit local path source", func(t *testing.T) {
		t.Parallel()
		setupCfg(t)

		src := t.TempDir()
		setupSourceSkill(t, src, "path-skill", "path-skill", "Path skill.")

		var buf bytes.Buffer
		c := listCmd
		c.SetArgs([]string{"-s", src})
		c.SetOut(&buf)
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("list with explicit path failed: %v", err)
		}
		if !strings.Contains(buf.String(), "path-skill") {
			t.Errorf("expected 'path-skill' in output, got %q", buf.String())
		}
	})

	t.Run("missing local with no default falls back to configured default", func(t *testing.T) {
		t.Parallel()
		// cfg has DefaultSource set; when cwd has no .agents and no --source,
		// it should fall back to cfg.DefaultSource
		dir := t.TempDir()
		cfg = &config.Config{DefaultSource: dir, DefaultRoot: "."}

		// Make dir have skills
		setupSourceSkill(t, dir, "fallback-skill", "fallback-skill", "Fallback skill.")

		cwd, _ := os.Getwd()
		_ = os.Chdir(t.TempDir()) // cwd with no .agents
		defer os.Chdir(cwd)

		var buf bytes.Buffer
		c := listCmd
		c.SetArgs([]string{})
		c.SetOut(&buf)
		c.SetErr(&bytes.Buffer{})

		if err := c.Execute(); err != nil {
			t.Fatalf("list with fallback source failed: %v", err)
		}
		if !strings.Contains(buf.String(), "fallback-skill") {
			t.Errorf("expected 'fallback-skill' in output, got %q", buf.String())
		}
	})

	t.Run("missing local with no default source errors", func(t *testing.T) {
		t.Parallel()
		cfg = &config.Config{DefaultSource: "", DefaultRoot: "."}

		dir := t.TempDir()
		cwd, _ := os.Getwd()
		_ = os.Chdir(dir) // no .agents here
		defer os.Chdir(cwd)

		c := listCmd
		c.SetArgs([]string{})
		c.SetOut(&bytes.Buffer{})
		c.SetErr(&bytes.Buffer{})

		err := c.Execute()
		if err == nil {
			t.Fatal("expected error when no local skills and no default source")
		}
	})
}
