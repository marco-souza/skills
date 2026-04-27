package skills

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CloneRepo clones a GitHub repo to a temp dir. Returns the dir and a cleanup func.
func CloneRepo(gh *GitHubSource) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "skills-repo-*")
	if err != nil {
		return "", nil, fmt.Errorf("creating temp dir: %w", err)
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }

	repoDir := filepath.Join(tmpDir, gh.Repo)
	cmd := exec.Command("git", "clone", "--depth", "1", gh.SSHURL, repoDir)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("cloning %s: %w", gh.SSHURL, err)
	}
	return repoDir, cleanup, nil
}
