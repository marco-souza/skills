package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// DefaultSkillsDir is the default hidden directory for agent skills.
	DefaultSkillsDir = ".agents"
	// SkillsSubDir is the subdirectory within the skills directory that contains skill definitions.
	SkillsSubDir = "skills"
	// SkillFileName is the standard filename for skill definition files.
	SkillFileName = "SKILL.md"
)

// Loader loads skills from a root project directory.
type Loader struct {
	// RootPath is the root directory of the project containing .agents/skills.
	RootPath string
}

// NewLoader returns a new Loader configured for the given root project path.
func NewLoader(rootPath string) *Loader {
	return &Loader{RootPath: rootPath}
}

// LoadAll discovers and loads all skills from the .agents/skills directory under RootPath.
// It skips entries that are not directories or lack a SKILL.md file, logging warnings for parse failures.
func (l *Loader) LoadAll() ([]*Skill, error) {
	skillsDir := filepath.Join(l.RootPath, DefaultSkillsDir, SkillsSubDir)

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

		skillPath := filepath.Join(skillsDir, entry.Name(), SkillFileName)
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

// ResolveToSkillsDir resolves an input path to the .agents/skills directory.
// If the input is relative, it is resolved against the current working directory.
// If the input already ends with .agents/skills, it is returned unchanged.
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
	if strings.HasSuffix(input, filepath.Join(DefaultSkillsDir, SkillsSubDir)) {
		return input
	}
	// Pointing at .agents — append skills
	if strings.HasSuffix(input, DefaultSkillsDir) {
		return filepath.Join(input, SkillsSubDir)
	}
	// Project root — append .agents/skills
	return filepath.Join(input, DefaultSkillsDir, SkillsSubDir)
}
