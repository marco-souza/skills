package skills

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ExecFunc is the type for the function used to spawn external commands.
type ExecFunc func(name string, args ...string) *exec.Cmd

// CloneRepo clones a GitHub repository to a temporary directory.
// execFn is the function used to spawn the git clone command; if nil, exec.Command is used.
// It returns the cloned repository path, a cleanup function to remove the temp directory, and any error.
func CloneRepo(gh *GitHubSource, execFn ExecFunc) (string, func(), error) {
	if execFn == nil {
		execFn = exec.Command
	}

	tmpDir, err := os.MkdirTemp("", "skills-repo-*")
	if err != nil {
		return "", nil, fmt.Errorf("creating temp dir: %w", err)
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }

	repoDir := filepath.Join(tmpDir, gh.Repo)
	cmd := execFn("git", "clone", "--depth", "1", gh.SSHURL, repoDir)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("cloning %s: %w", gh.SSHURL, err)
	}
	return repoDir, cleanup, nil
}
