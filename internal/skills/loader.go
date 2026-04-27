package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultSkillsDir = ".agents"
	skillsSubDir     = "skills"
	skillFileName    = "SKILL.md"
)

// Loader loads skills from a root project directory.
type Loader struct {
	RootPath string
}

// NewLoader creates a new Loader for the given root path.
func NewLoader(rootPath string) *Loader {
	return &Loader{RootPath: rootPath}
}

// LoadAll loads all skills from .agents/skills under RootPath.
func (l *Loader) LoadAll() ([]*Skill, error) {
	skillsDir := filepath.Join(l.RootPath, defaultSkillsDir, skillsSubDir)

	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("skills directory not found at %s", skillsDir)
		}
		return nil, fmt.Errorf("reading skills directory: %w", err)
	}

	var skills []*Skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillPath := filepath.Join(skillsDir, entry.Name(), skillFileName)
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			continue
		}

		skill := &Skill{}
		if err := skill.LoadFromPath(skillPath); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to load skill %s: %v\n", entry.Name(), err)
			continue
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

// ResolveToSkillsDir returns the .agents/skills directory for a given project root.
// Always appends .agents/skills unless the input already ends with it.
func ResolveToSkillsDir(input string) string {
	if input == "" {
		input = "."
	}
	if !filepath.IsAbs(input) {
		cwd, err := os.Getwd()
		if err == nil {
			input = filepath.Join(cwd, input)
		}
	}

	// Already pointing at .agents/skills
	if strings.HasSuffix(input, filepath.Join(defaultSkillsDir, skillsSubDir)) {
		return input
	}
	// Pointing at .agents — append skills
	if strings.HasSuffix(input, defaultSkillsDir) {
		return filepath.Join(input, skillsSubDir)
	}
	// Project root — append .agents/skills
	return filepath.Join(input, defaultSkillsDir, skillsSubDir)
}
