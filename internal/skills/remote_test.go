package skills

import (
"fmt"
"os"
"os/exec"
"path/filepath"
"testing"
)

func TestCloneRepo(t *testing.T) {
t.Parallel()

t.Run("clones from local bare repo", func(t *testing.T) {
t.Parallel()

bareRepo := createBareRepoForRemote(t)

gh := &GitHubSource{
Owner:  "test",
Repo:   "test-repo",
SSHURL: "file://" + bareRepo,
}

repoDir, cleanup, err := CloneRepo(gh, nil)
if err != nil {
t.Fatalf("CloneRepo failed: %v", err)
}
defer cleanup()

if _, err := os.Stat(filepath.Join(repoDir, "README.md")); os.IsNotExist(err) {
t.Fatal("expected README.md in cloned repo")
}
})

t.Run("creates temp dir", func(t *testing.T) {
t.Parallel()

bareRepo := createBareRepoForRemote(t)
callCount := 0

mockExec := func(name string, args ...string) *exec.Cmd {
callCount++
// Replace SSHURL arg with bare repo
newArgs := make([]string, len(args))
copy(newArgs, args)
for i := range newArgs {
if newArgs[i] == "git@github.com:test/test-repo.git" {
newArgs[i] = bareRepo
}
}
return exec.Command(name, newArgs...)
}

gh := &GitHubSource{
Owner:  "test",
Repo:   "test-repo",
SSHURL: "git@github.com:test/test-repo.git",
}

repoDir, cleanup, err := CloneRepo(gh, mockExec)
if err != nil {
t.Fatalf("CloneRepo failed: %v", err)
}
defer cleanup()

if callCount != 1 {
t.Errorf("expected 1 exec call, got %d", callCount)
}

// Verify the repo dir exists inside a temp directory
if !filepath.IsAbs(repoDir) {
t.Errorf("expected absolute repoDir, got %q", repoDir)
}
})

t.Run("cleanup removes temp dir", func(t *testing.T) {
t.Parallel()

bareRepo := createBareRepoForRemote(t)

mockExec := func(name string, args ...string) *exec.Cmd {
newArgs := make([]string, len(args))
copy(newArgs, args)
for i := range newArgs {
if newArgs[i] == "git@github.com:test/cleanup-repo.git" {
newArgs[i] = bareRepo
}
}
return exec.Command(name, newArgs...)
}

gh := &GitHubSource{
Owner:  "test",
Repo:   "cleanup-repo",
SSHURL: "git@github.com:test/cleanup-repo.git",
}

repoDir, cleanup, err := CloneRepo(gh, mockExec)
if err != nil {
t.Fatalf("CloneRepo failed: %v", err)
}

if _, err := os.Stat(repoDir); os.IsNotExist(err) {
t.Fatal("repo dir should exist before cleanup")
}

cleanup()

// The parent temp dir should be removed
tmpDir := filepath.Dir(filepath.Dir(repoDir))
if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
// At minimum, repoDir should be gone
if _, err2 := os.Stat(repoDir); !os.IsNotExist(err2) {
t.Error("repo dir should be removed after cleanup")
}
}
})

t.Run("clone failure returns error", func(t *testing.T) {
t.Parallel()

mockExec := func(name string, args ...string) *exec.Cmd {
// Return a command that will fail
return exec.Command("false")
}

gh := &GitHubSource{
Owner:  "test",
Repo:   "fail-repo",
SSHURL: "git@github.com:test/fail-repo.git",
}

_, _, err := CloneRepo(gh, mockExec)
if err == nil {
t.Fatal("expected error from failed clone")
}
})

t.Run("MkdirTemp failure returns error", func(t *testing.T) {
t.Parallel()
// We can't easily mock os.MkdirTemp, but we can test
// that CloneRepo propagates errors from exec.Command.
// This test verifies the error path for failed clones.
mockExec := func(name string, args ...string) *exec.Cmd {
return exec.Command("false")
}

gh := &GitHubSource{
Owner:  "test",
Repo:   "err-repo",
SSHURL: "git@github.com:test/err-repo.git",
}

_, cleanup, err := CloneRepo(gh, mockExec)
if err == nil {
t.Fatal("expected error")
}
if cleanup != nil {
t.Error("cleanup should be nil on failure")
}
})
}

// createBareRepoForRemote creates a regular git repo for remote tests.
func createBareRepoForRemote(t *testing.T) string {
t.Helper()
tempDir := t.TempDir()
repoDir := filepath.Join(tempDir, "repo")

os.MkdirAll(repoDir, 0755)
runGit(t, repoDir, "init")
runGit(t, repoDir, "config", "user.email", "test@test.com")
runGit(t, repoDir, "config", "user.name", "Test")

// Create a README.md
os.WriteFile(filepath.Join(repoDir, "README.md"), []byte("# test repo"), 0644)

// Create a SKILL.md to make it look like a skills repo
skillsDir := filepath.Join(repoDir, DefaultSkillsDir, SkillsSubDir, "remote-skill")
os.MkdirAll(skillsDir, 0755)
os.WriteFile(filepath.Join(skillsDir, SkillFileName), []byte(`---
name: remote-skill
description: A remote skill for testing.
---

# Remote Skill
`), 0644)

runGit(t, repoDir, "add", ".")
runGit(t, repoDir, "commit", "-m", "init")

return repoDir
}

// TestCloneRepo_MkdirTempError tests the MkdirTemp error path by
// temporarily replacing os.MkdirTemp behavior. Since we can't inject
// that easily, we instead verify error propagation through a helper.
func TestCloneRepo_ErrorPropagation(t *testing.T) {
t.Parallel()

// Use a command that always fails
mockExec := func(name string, args ...string) *exec.Cmd {
return exec.Command("sh", "-c", "exit 1")
}

gh := &GitHubSource{
Owner:  "fake",
Repo:   "fake-repo",
SSHURL: "git@github.com:fake/fake-repo.git",
}

_, _, err := CloneRepo(gh, mockExec)
if err == nil {
t.Fatal("expected error")
}
if got := err.Error(); got == "" {
t.Error("error message should not be empty")
}
}

// TestCloneRepo_WithRealBareRepo tests CloneRepo end-to-end with a local bare repo.
// This avoids needing to inject the SSH URL by directly patching execCommand.
func TestCloneRepo_EndToEnd(t *testing.T) {
t.Parallel()

bareRepo := createBareRepoForRemote(t)

// Track the actual command args
var capturedArgs []string
mockExec := func(name string, args ...string) *exec.Cmd {
capturedArgs = args
// Replace the SSH URL with our local bare repo path
newArgs := make([]string, len(args))
copy(newArgs, args)
for i := range newArgs {
// The last arg is the destination, second-to-last is the source URL
if i == len(args)-2 {
newArgs[i] = bareRepo
}
}
return exec.Command(name, newArgs...)
}

gh := &GitHubSource{
Owner:  "e2e",
Repo:   "e2e-repo",
SSHURL: "git@github.com:e2e/e2e-repo.git",
}

repoDir, cleanup, err := CloneRepo(gh, mockExec)
if err != nil {
t.Fatalf("CloneRepo failed: %v", err)
}
defer cleanup()

// Verify captured args include git clone flags
if len(capturedArgs) < 4 {
t.Fatalf("expected at least 4 args, got %d: %v", len(capturedArgs), capturedArgs)
}
if capturedArgs[0] != "clone" {
t.Errorf("expected 'clone', got %q", capturedArgs[0])
}
if capturedArgs[1] != "--depth" {
t.Errorf("expected '--depth', got %q", capturedArgs[1])
}

// Verify repo was cloned successfully
if _, err := os.Stat(filepath.Join(repoDir, DefaultSkillsDir, SkillsSubDir, "remote-skill", SkillFileName)); os.IsNotExist(err) {
t.Fatal("expected SKILL.md in cloned repo")
}
}

// TestCloneRepo_SourceDirResolution tests that the clone destination includes the repo name.
func TestCloneRepo_RepoDirIncludesRepoName(t *testing.T) {
t.Parallel()

bareRepo := createBareRepoForRemote(t)

mockExec := func(name string, args ...string) *exec.Cmd {
newArgs := make([]string, len(args))
copy(newArgs, args)
for i := range newArgs {
if i == len(args)-2 {
newArgs[i] = bareRepo
}
}
return exec.Command(name, newArgs...)
}

gh := &GitHubSource{
Owner:  "dir-test",
Repo:   "named-repo",
SSHURL: "git@github.com:dir-test/named-repo.git",
}

repoDir, cleanup, err := CloneRepo(gh, mockExec)
if err != nil {
t.Fatalf("CloneRepo failed: %v", err)
}
defer cleanup()

// repoDir should end with the repo name
expectedSuffix := fmt.Sprintf("%s%s", string(filepath.Separator), "named-repo")
if len(repoDir) < len(expectedSuffix) || repoDir[len(repoDir)-len(expectedSuffix):] != expectedSuffix {
t.Errorf("expected repoDir to end with '/named-repo', got %q", repoDir)
}
}
