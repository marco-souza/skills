package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFrontmatter_Valid(t *testing.T) {
	content := `---
name: test-skill
description: >
  A test skill.
  Use when: testing.
  Do NOT use when: production.
---

# Test Skill

Body content here.
`
	fm, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm["name"] != "test-skill" {
		t.Errorf("expected name 'test-skill', got %v", fm["name"])
	}
	if body != "# Test Skill\n\nBody content here." {
		t.Errorf("unexpected body: %q", body)
	}
}

func TestParseFrontmatter_Missing(t *testing.T) {
	_, _, err := ParseFrontmatter("# No frontmatter")
	if err == nil {
		t.Fatal("expected error for missing frontmatter")
	}
}

func TestParseFrontmatter_InvalidYAML(t *testing.T) {
	t.Parallel()
	content := `---
name: [broken
---
body`
	_, _, err := ParseFrontmatter(content)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestParseFrontmatter_MissingClosing(t *testing.T) {
	t.Parallel()
	content := `---
name: test
description: no closing dashes`
	_, _, err := ParseFrontmatter(content)
	if err == nil {
		t.Fatal("expected error for missing closing ---")
	}
}

func TestParseFrontmatter_EmptyBody(t *testing.T) {
	t.Parallel()
	content := `---
name: no-body
description: desc
---`
	fm, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm["name"] != "no-body" {
		t.Errorf("expected name 'no-body', got %v", fm["name"])
	}
	if body != "" {
		t.Errorf("expected empty body, got %q", body)
	}
}

func TestSkill_LoadFromPath(t *testing.T) {
	skill := &Skill{}
	err := skill.LoadFromPath("../../testdata/valid-skill/SKILL.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skill.Name != "valid-skill" {
		t.Errorf("expected name 'valid-skill', got %q", skill.Name)
	}
	if skill.Description == "" {
		t.Error("expected non-empty description")
	}
}

func TestSkill_Validate_Valid(t *testing.T) {
	skill := &Skill{
		Name:        "valid-skill",
		Description: "A test. Use when: testing. Do NOT use when: production.",
	}
	if err := skill.Validate(); err != nil {
		t.Errorf("expected valid skill, got error: %v", err)
	}
}

func TestSkill_Validate_BadName(t *testing.T) {
	skill := &Skill{
		Name:        "BAD_NAME",
		Description: "A test. Use when: testing. Do NOT use when: production.",
	}
	if err := skill.Validate(); err == nil {
		t.Fatal("expected validation error for uppercase name")
	}
}

func TestSkill_Validate_DescriptionTooLong(t *testing.T) {
	t.Parallel()
	s := &Skill{
		Name:        "test",
		Description: strings.Repeat("x", 1025),
	}
	err := s.Validate()
	if err == nil {
		t.Fatal("expected validation error for description too long")
	}
	if !strings.Contains(err.Error(), "1024") {
		t.Errorf("expected error mentioning 1024, got %v", err)
	}
}

func TestSkill_Validate_EmptyName(t *testing.T) {
	t.Parallel()
	s := &Skill{Name: "", Description: "A description."}
	err := s.Validate()
	if err == nil {
		t.Fatal("expected validation error for empty name")
	}
	if !strings.Contains(err.Error(), "name is required") {
		t.Errorf("expected 'name is required', got %v", err)
	}
}

func TestSkill_Validate_EmptyDescription(t *testing.T) {
	t.Parallel()
	s := &Skill{Name: "test-skill", Description: ""}
	err := s.Validate()
	if err == nil {
		t.Fatal("expected validation error for empty description")
	}
	if !strings.Contains(err.Error(), "description is required") {
		t.Errorf("expected 'description is required', got %v", err)
	}
}

func TestSkill_Validate_SingleCharName(t *testing.T) {
	t.Parallel()
	s := &Skill{Name: "a", Description: "A test."}
	if err := s.Validate(); err != nil {
		t.Errorf("single char name should be valid: %v", err)
	}
}

func TestSkill_LoadFromPath_MissingFile(t *testing.T) {
	t.Parallel()
	s := &Skill{}
	if err := s.LoadFromPath("/nonexistent/SKILL.md"); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSkill_LoadFromPath_InvalidYAML(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	skillPath := filepath.Join(dir, SkillFileName)
	if err := os.WriteFile(skillPath, []byte("---\ninvalid: [yaml\n---\nbody"), 0o644); err != nil {
		t.Fatalf("writing invalid skill file: %v", err)
	}

	s := &Skill{}
	if err := s.LoadFromPath(skillPath); err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestSkill_LoadFromPath_FallbackToDirName(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "fallback-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("creating fallback skill dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, SkillFileName), []byte("---\ndescription: desc\n---\nbody"), 0o644); err != nil {
		t.Fatalf("writing fallback skill file: %v", err)
	}

	s := &Skill{}
	if err := s.LoadFromPath(filepath.Join(skillDir, SkillFileName)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != "fallback-skill" {
		t.Errorf("expected name 'fallback-skill' from dir, got %q", s.Name)
	}
}

func TestSkill_LoadFromPath_NameAndMetadata(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	content := `---
name: full-skill
description: A full skill.
metadata:
  runtime: python3
  scripts:
    - scripts/run.sh
---

# Full Skill
Body text.
	`
	skillPath := filepath.Join(dir, SkillFileName)
	if err := os.WriteFile(skillPath, []byte(content), 0o644); err != nil {
		t.Fatalf("writing skill file: %v", err)
	}

	s := &Skill{}
	if err := s.LoadFromPath(skillPath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != "full-skill" {
		t.Errorf("expected name 'full-skill', got %q", s.Name)
	}
	if s.Description != "A full skill." {
		t.Errorf("expected description 'A full skill.', got %q", s.Description)
	}
	if s.Runtime() != "python3" {
		t.Errorf("expected runtime 'python3', got %q", s.Runtime())
	}
	if s.Path != skillPath {
		t.Errorf("expected path %q, got %q", skillPath, s.Path)
	}
	if s.Content != "# Full Skill\nBody text." {
		t.Errorf("unexpected content: %q", s.Content)
	}
}

func TestValidationError_Error(t *testing.T) {
	t.Parallel()
	e := &ValidationError{Errors: []string{"err1", "err2"}}
	msg := e.Error()
	if !strings.Contains(msg, "err1") || !strings.Contains(msg, "err2") || !strings.Contains(msg, "validation failed") {
		t.Errorf("unexpected error message: %q", msg)
	}
}

func TestSkill_Scripts(t *testing.T) {
	t.Parallel()

	t.Run("valid list", func(t *testing.T) {
		t.Parallel()
		s := &Skill{Metadata: map[string]any{"scripts": []any{"scripts/setup.sh", "scripts/cleanup.sh"}}}
		scripts := s.Scripts()
		if len(scripts) != 2 || scripts[0] != "scripts/setup.sh" || scripts[1] != "scripts/cleanup.sh" {
			t.Errorf("unexpected scripts: %v", scripts)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		t.Parallel()
		s := &Skill{Metadata: map[string]any{"other": "value"}}
		if len(s.Scripts()) != 0 {
			t.Errorf("expected 0 scripts, got %d", len(s.Scripts()))
		}
	})

	t.Run("nil metadata", func(t *testing.T) {
		t.Parallel()
		if len((&Skill{}).Scripts()) != 0 {
			t.Error("expected empty scripts for nil metadata")
		}
	})

	t.Run("mixed types", func(t *testing.T) {
		t.Parallel()
		s := &Skill{Metadata: map[string]any{"scripts": []any{"a.sh", 42, nil, "b.sh"}}}
		scripts := s.Scripts()
		if len(scripts) != 2 || scripts[0] != "a.sh" || scripts[1] != "b.sh" {
			t.Errorf("unexpected scripts: %v", scripts)
		}
	})
}

func TestSkill_Runtime(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		s := &Skill{Metadata: map[string]any{"runtime": "python3"}}
		if s.Runtime() != "python3" {
			t.Errorf("expected 'python3', got %q", s.Runtime())
		}
	})

	t.Run("nil metadata", func(t *testing.T) {
		t.Parallel()
		if (&Skill{}).Runtime() != "" {
			t.Error("expected empty for nil metadata")
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		t.Parallel()
		s := &Skill{Metadata: map[string]any{"runtime": 42}}
		if s.Runtime() != "" {
			t.Errorf("expected empty for wrong type, got %q", s.Runtime())
		}
	})
}

func TestSkill_Dependencies(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		s := &Skill{Metadata: map[string]any{
			"dependencies": map[string]any{"skills": []any{"git-commit-formatter", "pr-review"}},
		}}
		deps := s.Dependencies()
		if len(deps) != 2 || deps[0] != "git-commit-formatter" || deps[1] != "pr-review" {
			t.Errorf("unexpected deps: %v", deps)
		}
	})

	t.Run("nil metadata", func(t *testing.T) {
		t.Parallel()
		if (&Skill{}).Dependencies() != nil {
			t.Error("expected nil deps for nil metadata")
		}
	})

	t.Run("mixed types", func(t *testing.T) {
		t.Parallel()
		s := &Skill{Metadata: map[string]any{
			"dependencies": map[string]any{"skills": []any{"valid-skill", 42, nil, "another-skill"}},
		}}
		deps := s.Dependencies()
		if len(deps) != 2 || deps[0] != "valid-skill" || deps[1] != "another-skill" {
			t.Errorf("unexpected deps: %v", deps)
		}
	})
}

func TestResolveToSkillsDir(t *testing.T) {
	result := ResolveToSkillsDir("")
	if result == "" {
		t.Fatal("expected non-empty result for empty input")
	}
	if !filepath.IsAbs(result) {
		t.Errorf("expected absolute path, got %q", result)
	}
	if !strings.HasSuffix(result, filepath.Join(".agents", "skills")) {
		t.Errorf("expected path ending in .agents/skills, got %q", result)
	}
}
