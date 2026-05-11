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

// Source is a sum type representing either a local filesystem path or a GitHub repo.
// Exactly one of the fields is non-nil.
type Source struct {
	// Local is set when the source is a local filesystem path.
	Local *LocalSource
	// GitHub is set when the source is a GitHub repository.
	GitHub *GitHubSource
}

// LocalSource represents a local filesystem path.
type LocalSource struct {
	// Path is the absolute filesystem path to the skills directory.
	Path string
}

// GitHubSource represents a GitHub repository.
type GitHubSource struct {
	// Owner is the GitHub account or organization name.
	Owner string
	// Repo is the repository name.
	Repo string
	// URL is the HTTPS clone URL.
	URL string
	// SSHURL is the SSH clone URL.
	SSHURL string
}

func ghSource(owner, repo string) *GitHubSource {
	return &GitHubSource{
		Owner:  owner,
		Repo:   repo,
		URL:    fmt.Sprintf("https://github.com/%s/%s", owner, repo),
		SSHURL: fmt.Sprintf("git@github.com:%s/%s.git", owner, repo),
	}
}

// ResolvePath resolves a user-provided input string to a Source.
// It interprets the input as a GitHub owner/repo shorthand, a GitHub URL, or a local filesystem path.
func ResolvePath(input string) (Source, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return Source{}, fmt.Errorf("input cannot be empty")
	}

	if gh := resolveGitHub(input); gh != nil {
		return Source{GitHub: gh}, nil
	}

	cwd, _ := os.Getwd()
	path := input
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	return Source{Local: &LocalSource{Path: path}}, nil
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

// ResolveSourceDir resolves a source string to an actual directory path
// containing skills, with an optional cleanup function.
//
// Resolution order:
//  1. If source is empty, check for local .agents/skills at cwd.
//     If found, return "." with no cleanup.
//  2. If source is empty and local not found, fall back to defaultSource.
//  3. If source is non-empty, resolve it via ResolvePath.
//
// If the resolved source is a GitHub repo, it is cloned to a temp dir
// and the cleanup function removes it when called.
// execFn is the function used to spawn the git clone command; if nil, exec.Command is used.
func ResolveSourceDir(source string, defaultSource string, execFn ExecFunc) (string, func(), error) {
	// No explicit source: prefer local .agents/skills
	if source == "" {
		localDir := ResolveToSkillsDir(".")
		if _, err := os.Stat(localDir); err == nil {
			return ".", nil, nil
		}
		// Fall back to defaultSource
		if defaultSource == "" {
			return "", nil, fmt.Errorf("no local skills found and no --source specified")
		}
		source = defaultSource
	}

	// Resolve the source: local path or GitHub repo
	src, err := ResolvePath(source)
	if err != nil {
		return "", nil, fmt.Errorf("resolving source %q: %w", source, err)
	}

	// Sum-type dispatch: exactly one field is non-nil
	if src.Local != nil {
		return src.Local.Path, nil, nil
	}
	if src.GitHub != nil {
		tmpDir, cleanup, err := CloneRepo(src.GitHub, execFn)
		if err != nil {
			return "", nil, err
		}
		return tmpDir, cleanup, nil
	}

	return "", nil, fmt.Errorf("--source requires a GitHub repo (owner/repo) or local path, got %s", source)
}
