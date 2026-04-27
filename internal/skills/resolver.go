package skills

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var githubShorthandRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+/[a-zA-Z0-9_.-]+$`)

// Source represents where skills come from.
type Source interface {
	Type() string
	String() string
}

// LocalSource represents a local filesystem path.
type LocalSource struct{ Path string }

func (l *LocalSource) Type() string { return "local" }
func (l *LocalSource) String() string { return l.Path }

// GitHubSource represents a GitHub repository.
type GitHubSource struct {
	Owner  string
	Repo   string
	URL    string // HTTPS
	SSHURL string // SSH
}

func (g *GitHubSource) Type() string { return "github" }
func (g *GitHubSource) String() string {
	return fmt.Sprintf("github.com/%s/%s", g.Owner, g.Repo)
}

// ResolvePath resolves input to a Source (local path or GitHub repo).
func ResolvePath(input string) (Source, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("input cannot be empty")
	}

	if gh := resolveGitHub(input); gh != nil {
		return gh, nil
	}

	cwd, _ := os.Getwd()
	path := input
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	return &LocalSource{Path: path}, nil
}

func resolveGitHub(input string) *GitHubSource {
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
		URL:    fmt.Sprintf("https://github.com/%s/%s", owner, repo),
		SSHURL: fmt.Sprintf("git@github.com:%s/%s.git", owner, repo),
	}
}
