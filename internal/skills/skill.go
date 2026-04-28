package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Skill represents an AI agent skill definition.
type Skill struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Tags        []string               `yaml:"tags,omitempty"`
	Category    string                 `yaml:"category,omitempty"`
	Author      string                 `yaml:"author,omitempty"`
	Version     string                 `yaml:"version,omitempty"`
	Metadata    map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Path        string                 `json:"-"`
	Content     string                 `json:"-"`
}

var nameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,62}[a-z0-9]$|^[a-z0-9]$`)

// Validate checks that the skill meets format requirements.
func (s *Skill) Validate() error {
	var errs []string

	if s.Name == "" {
		errs = append(errs, "name is required")
	} else if !nameRegex.MatchString(s.Name) {
		errs = append(errs, "name must be 1-64 lowercase letters, numbers, and hyphens")
	}

	if s.Description == "" {
		errs = append(errs, "description is required")
	} else if len(s.Description) > 1024 {
		errs = append(errs, fmt.Sprintf("description is %d chars (max 1024)", len(s.Description)))
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}
	return nil
}

// ValidationError holds multiple validation errors.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return "validation failed:\n  - " + strings.Join(e.Errors, "\n  - ")
}

// LoadFromPath loads a skill from a SKILL.md file.
func (s *Skill) LoadFromPath(path string) error {
	s.Path = path

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading SKILL.md: %w", err)
	}

	frontmatter, body, err := ParseFrontmatter(string(data))
	if err != nil {
		return fmt.Errorf("parsing frontmatter: %w", err)
	}

	dirName := filepath.Base(filepath.Dir(path))

	if name, ok := frontmatter["name"].(string); ok && name != "" {
		s.Name = name
	} else {
		s.Name = dirName
	}

	if desc, ok := frontmatter["description"].(string); ok {
		s.Description = desc
	}
	if tags, ok := frontmatter["tags"].([]interface{}); ok {
		for _, t := range tags {
			if ts, ok := t.(string); ok {
				s.Tags = append(s.Tags, ts)
			}
		}
	}
	if cat, ok := frontmatter["category"].(string); ok {
		s.Category = cat
	}
	if author, ok := frontmatter["author"].(string); ok {
		s.Author = author
	}
	if version, ok := frontmatter["version"].(string); ok {
		s.Version = version
	}
	if metadata, ok := frontmatter["metadata"].(map[string]interface{}); ok {
		s.Metadata = metadata
	}

	s.Content = strings.TrimSpace(body)
	return nil
}

// Scripts returns the list of script dependencies from metadata.scripts.
// Returns an empty slice if metadata.scripts is absent or not a list of strings.
func (s *Skill) Scripts() []string {
	if s.Metadata == nil {
		return []string{}
	}
	scriptsRaw, ok := s.Metadata["scripts"].([]interface{})
	if !ok {
		return []string{}
	}
	scripts := make([]string, 0, len(scriptsRaw))
	for _, item := range scriptsRaw {
		if script, ok := item.(string); ok {
			scripts = append(scripts, script)
		}
	}
	return scripts
}

// Runtime returns the runtime identifier from metadata.runtime.
// Returns an empty string if metadata.runtime is absent or not a string.
func (s *Skill) Runtime() string {
	if s.Metadata == nil {
		return ""
	}
	if runtime, ok := s.Metadata["runtime"].(string); ok {
		return runtime
	}
	return ""
}
