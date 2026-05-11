// Package skills provides types and functions for loading, validating, resolving, and installing AI agent skills.
package skills

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Skill represents an AI agent skill definition loaded from a SKILL.md file.
type Skill struct {
	// Name is the skill identifier (lowercase, hyphenated).
	Name string `yaml:"name"`
	// Description explains what the skill does.
	Description string `yaml:"description"`
	// Tags are optional labels for categorising the skill.
	Tags []string `yaml:"tags,omitempty"`
	// Category is an optional high-level grouping.
	Category string `yaml:"category,omitempty"`
	// Author is the skill's creator.
	Author string `yaml:"author,omitempty"`
	// Version is the skill's semantic version.
	Version string `yaml:"version,omitempty"`
	// Metadata holds arbitrary extra data from the frontmatter.
	Metadata map[string]any `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	// Path is the filesystem path to the SKILL.md file.
	Path string `json:"-"`
	// Content is the markdown body after the frontmatter.
	Content string `json:"-"`
}

var nameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,62}[a-z0-9]$|^[a-z0-9]$`)

// Validate checks that the skill meets the required format constraints.
// It returns a ValidationError if the name or description is invalid.
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

// ValidationError holds a list of validation error messages.
type ValidationError struct {
	// Errors is the list of individual validation failure messages.
	Errors []string
}

// Error returns a formatted string of all validation error messages.
func (e *ValidationError) Error() string {
	return "validation failed:\n  - " + strings.Join(e.Errors, "\n  - ")
}

// LoadFromPath reads and parses a SKILL.md file into the Skill struct.
// It extracts YAML frontmatter fields and stores the remaining markdown as Content.
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
	if tags, ok := frontmatter["tags"].([]any); ok {
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
	if metadata, ok := frontmatter["metadata"].(map[string]any); ok {
		s.Metadata = metadata
	}

	s.Content = strings.TrimSpace(body)
	return nil
}

// Scripts returns the list of script file paths declared in metadata.scripts.
// It returns an empty slice if the field is absent or malformed.
func (s *Skill) Scripts() []string {
	if s.Metadata == nil {
		return []string{}
	}
	scriptsRaw, ok := s.Metadata["scripts"].([]any)
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
// It returns an empty string if the field is absent or not a string.
func (s *Skill) Runtime() string {
	if s.Metadata == nil {
		return ""
	}
	if runtime, ok := s.Metadata["runtime"].(string); ok {
		return runtime
	}
	return ""
}

// Dependencies returns the list of skill names from metadata.dependencies.skills.
// It returns nil if the dependencies section is absent or malformed.
func (s *Skill) Dependencies() []string {
	if s.Metadata == nil {
		return nil
	}
	depsRaw, ok := s.Metadata["dependencies"].(map[string]any)
	if !ok {
		return nil
	}
	skillsRaw, ok := depsRaw["skills"].([]any)
	if !ok {
		return nil
	}
	deps := make([]string, 0, len(skillsRaw))
	for _, item := range skillsRaw {
		if name, ok := item.(string); ok {
			deps = append(deps, name)
		}
	}
	return deps
}
