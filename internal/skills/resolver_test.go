package skills

import (
"os"
"path/filepath"
"testing"
)

func TestResolveSourceDir(t *testing.T) {
t.Parallel()

t.Run("empty source with local .agents/skills found", func(t *testing.T) {
t.Parallel()
dir := t.TempDir()
os.MkdirAll(filepath.Join(dir, DefaultSkillsDir, SkillsSubDir), 0755)

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

t.Run("empty source with no local and no default errors", func(t *testing.T) {
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

t.Run("empty source falls back to local default source", func(t *testing.T) {
t.Parallel()
dir := t.TempDir()
cwd, _ := os.Getwd()
_ = os.Chdir(dir)
defer os.Chdir(cwd)

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

t.Run("absolute local source", func(t *testing.T) {
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

t.Run("relative local source", func(t *testing.T) {
t.Parallel()
dir := t.TempDir()
cwd, _ := os.Getwd()
_ = os.Chdir(dir)
defer os.Chdir(cwd)

os.MkdirAll("subdir", 0755)
result, cleanup, err := ResolveSourceDir("subdir", "")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if cleanup != nil {
t.Fatal("cleanup should be nil for local source")
}
expected := filepath.Clean(filepath.Join(dir, "subdir"))
if filepath.Clean(result) != expected {
t.Errorf("expected %q, got %q", expected, result)
}
})
}

func TestResolveGitHub(t *testing.T) {
t.Parallel()

t.Run("shorthand owner/repo", func(t *testing.T) {
t.Parallel()
gh := ResolveGitHub("user/repo")
if gh == nil {
t.Fatal("expected GitHub source")
}
if gh.Owner != "user" || gh.Repo != "repo" {
t.Errorf("expected user/repo, got %s/%s", gh.Owner, gh.Repo)
}
})

t.Run("https URL", func(t *testing.T) {
t.Parallel()
gh := ResolveGitHub("https://github.com/user/repo")
if gh == nil {
t.Fatal("expected GitHub source")
}
if gh.Owner != "user" || gh.Repo != "repo" {
t.Errorf("expected user/repo, got %s/%s", gh.Owner, gh.Repo)
}
})

t.Run("absolute path returns nil", func(t *testing.T) {
t.Parallel()
if ResolveGitHub("/absolute/path") != nil {
t.Error("expected nil for absolute path")
}
})

t.Run("local relative path returns nil", func(t *testing.T) {
t.Parallel()
if ResolveGitHub("./some/path") != nil {
t.Error("expected nil for relative path with ./")
}
})
}
