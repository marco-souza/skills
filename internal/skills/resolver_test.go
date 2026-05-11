package skills

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestResolveSourceDir(t *testing.T) {
	t.Parallel()

	t.Run("empty source with local .agents/skills found", func(t *testing.T) {
		t.Parallel()
		// Create a temp dir with .agents/skills
		dir := t.TempDir()
		skillsDir := filepath.Join(dir, DefaultSkillsDir, SkillsSubDir)
		os.MkdirAll(skillsDir, 0755)

		// Change to the temp dir
		cwd, _ := os.Getwd()
		_ = os.Chdir(dir)
		defer os.Chdir(cwd)

		result, cleanup, err := ResolveSourceDir("", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cleanup != nil {
			t.Fatal("cleanup should be nil for local source")
		}
		if result != "." {
			t.Errorf("expected '.', got %q", result)
		}
	})

	t.Run("empty source with no local and no default", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		cwd, _ := os.Getwd()
		_ = os.Chdir(dir)
		defer os.Chdir(cwd)

		_, _, err := ResolveSourceDir("", "")
		if err == nil {
			t.Fatal("expected error when no local skills and no default source")
		}
	})

	t.Run("empty source falls back to default source", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		cwd, _ := os.Getwd()
		_ = os.Chdir(dir)
		defer os.Chdir(cwd)

		// Use a local path as default source
		defaultSrc := "/some/local/path"
		result, cleanup, err := ResolveSourceDir("", defaultSrc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cleanup != nil {
			t.Fatal("cleanup should be nil for local default source")
		}
		if result != defaultSrc {
			t.Errorf("expected %q, got %q", defaultSrc, result)
		}
	})

	t.Run("non-empty local source", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		result, cleanup, err := ResolveSourceDir(dir, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cleanup != nil {
			t.Fatal("cleanup should be nil for local source")
		}
		if result != dir {
			t.Errorf("expected %q, got %q", dir, result)
		}
	})

	t.Run("non-empty relative source", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		cwd, _ := os.Getwd()
		_ = os.Chdir(dir)
		defer os.Chdir(cwd)

		// Create a subdirectory
		subDir := "subdir"
		os.MkdirAll(subDir, 0755)

		result, cleanup, err := ResolveSourceDir(subDir, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cleanup != nil {
			t.Fatal("cleanup should be nil for local source")
		}
		// Use filepath.Clean for consistent path comparison (macOS /private/ prefix)
		expected := filepath.Clean(filepath.Join(dir, subDir))
		resultClean := filepath.Clean(result)
		if resultClean != expected {
			t.Errorf("expected %q, got %q", expected, resultClean)
		}
	})

	t.Run("github source clones repo", func(t *testing.T) {
		t.Parallel()
		origExec := execCommand
		defer func() { execCommand = origExec }()

		bareRepo := createBareRepoForResolver(t)

		execCommand = func(name string, args ...string) *exec.Cmd {
			newArgs := make([]string, len(args))
			copy(newArgs, args)
			for i := range newArgs {
				if i == len(args)-2 {
					newArgs[i] = bareRepo
				}
			}
			return exec.Command(name, newArgs...)
		}

		// Use a GitHub shorthand as source
		result, cleanup, err := ResolveSourceDir("resolver-test/repo", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cleanup == nil {
			t.Fatal("cleanup should not be nil for GitHub source")
		}
		defer cleanup()

		// Result should be an absolute temp directory path
		if !filepath.IsAbs(result) {
			t.Errorf("expected absolute path, got %q", result)
		}

		// The cloned repo should contain our test skill
		skillPath := filepath.Join(result, "repo", DefaultSkillsDir, SkillsSubDir, "resolver-skill", SkillFileName)
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			t.Fatalf("expected SKILL.md at %s", skillPath)
		}
	})

	t.Run("github source with clone failure", func(t *testing.T) {
		t.Parallel()
		origExec := execCommand
		defer func() { execCommand = origExec }()

		execCommand = func(name string, args ...string) *exec.Cmd {
			return exec.Command("false")
		}

		_, _, err := ResolveSourceDir("fail-test/fail-repo", "")
		if err == nil {
			t.Fatal("expected error from failed clone")
		}
	})

	t.Run("invalid empty source string", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		cwd, _ := os.Getwd()
		_ = os.Chdir(dir)
		defer os.Chdir(cwd)

		// Empty input with no local and no default should error
		_, _, err := ResolveSourceDir("", "")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func createBareRepoForResolver(t *testing.T) string {
	t.Helper()
	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "repo")

	os.MkdirAll(repoDir, 0755)
	runGit(t, repoDir, "init")
	runGit(t, repoDir, "config", "user.email", "test@test.com")
	runGit(t, repoDir, "config", "user.name", "Test")

	skillsDir := filepath.Join(repoDir, DefaultSkillsDir, SkillsSubDir, "resolver-skill")
	os.MkdirAll(skillsDir, 0755)
	os.WriteFile(filepath.Join(skillsDir, SkillFileName), []byte(`---
name: resolver-skill
description: A resolver test skill.
---

# Resolver Skill
`), 0644)

	runGit(t, repoDir, "add", ".")
	runGit(t, repoDir, "commit", "-m", "init")

	return repoDir
}
