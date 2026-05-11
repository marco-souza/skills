package skills

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// ExecFunc is the type for the function used to spawn external commands.
type ExecFunc func(name string, args ...string) *exec.Cmd

// GitHubSource represents a GitHub repository.
type GitHubSource struct {
	// Owner is the GitHub account or organization name.
	Owner string
	// Repo is the repository name.
	Repo string
	// SSHURL is the SSH clone URL.
	SSHURL string
}

var githubShorthandRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+/[a-zA-Z0-9_.-]+$`)

// ResolveGitHub parses a user-provided string into a GitHubSource.
// It accepts owner/repo shorthands and https://github.com/... URLs.
// Returns nil if the input is not a GitHub reference.
func ResolveGitHub(input string) *GitHubSource {
	if filepath.IsAbs(input) {
		return nil
	}

	// Shorthand: owner/repo
	if githubShorthandRe.MatchString(input) {
		parts := strings.SplitN(input, "/", 2)
		return ghSource(parts[0], strings.TrimSuffix(parts[1], ".git"))
	}

	// URL: https://github.com/... or github.com/...
	raw := input
	if strings.HasPrefix(raw, "github.com/") {
		raw = "https://" + raw
	}
	if !strings.HasPrefix(raw, "https://") {
		return nil
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host != "github.com" {
		return nil
	}
	path := strings.TrimPrefix(u.Path, "/")
	parts := strings.SplitN(path, "/", 3)
	if len(parts) < 2 {
		return nil
	}
	return ghSource(parts[0], strings.TrimSuffix(parts[1], ".git"))
}

func ghSource(owner, repo string) *GitHubSource {
	return &GitHubSource{
		Owner:  owner,
		Repo:   repo,
		SSHURL: fmt.Sprintf("git@github.com:%s/%s.git", owner, repo),
	}
}

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
