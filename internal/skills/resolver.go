package skills

import (
	"fmt"
	"os"
	"path/filepath"
)

// ResolveSourceDir resolves a source string to an actual directory path
// containing skills, with an optional cleanup function.
//
// Resolution order:
//  1. If source is empty, check for local .agents/skills at cwd.
//     If found, return "." with no cleanup.
//  2. If source is empty and local not found, fall back to defaultSource.
//  3. If source is non-empty, check if it's a GitHub reference and clone,
//     or treat as a local path.
func ResolveSourceDir(source string, defaultSource string) (string, func(), error) {
	// No explicit source: prefer local .agents/skills
	if source == "" {
		localDir := ResolveToSkillsDir(".")
		if _, err := os.Stat(localDir); err == nil {
			return ".", nil, nil
		}
		if defaultSource == "" {
			return "", nil, fmt.Errorf("no local skills found and no --source specified")
		}
		source = defaultSource
	}

	// GitHub source
	if gh := ResolveGitHub(source); gh != nil {
		repoDir, cleanup, err := CloneRepo(gh, nil)
		if err != nil {
			return "", nil, err
		}
		return repoDir, cleanup, nil
	}

	// Local path
	if !filepath.IsAbs(source) {
		cwd, _ := os.Getwd()
		source = filepath.Join(cwd, source)
	}
	return source, nil, nil
}
