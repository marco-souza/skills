package skills

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// execCommand is the function used to spawn git clone commands.
// It is a package-level variable so tests can inject a mock.
var execCommand = exec.Command

// CloneRepo clones a GitHub repository to a temporary directory.
// It returns the cloned repository path, a cleanup function to remove the temp directory, and any error.
func CloneRepo(gh *GitHubSource) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "skills-repo-*")
	if err != nil {
		return "", nil, fmt.Errorf("creating temp dir: %w", err)
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }

	repoDir := filepath.Join(tmpDir, gh.Repo)
	cmd := execCommand("git", "clone", "--depth", "1", gh.SSHURL, repoDir)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("cloning %s: %w", gh.SSHURL, err)
	}
	return repoDir, cleanup, nil
}
