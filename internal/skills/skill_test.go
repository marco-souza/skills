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
		Path:        "testdata/valid-skill/SKILL.md",
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
	err := skill.Validate()
	if err == nil {
		t.Fatal("expected validation error for uppercase name")
	}
}

func TestResolvePath_Local(t *testing.T) {
	src, err := ResolvePath("/absolute/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Local == nil {
		t.Fatalf("expected Local source, got %+v", src)
	}
	if src.Local.Path != "/absolute/path" {
		t.Errorf("expected '/absolute/path', got %q", src.Local.Path)
	}
}

func TestResolvePath_Relative(t *testing.T) {
	src, err := ResolvePath("./some/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Local == nil {
		t.Fatalf("expected Local source, got %+v", src)
	}
	cwd, _ := os.Getwd()
	expected := filepath.Join(cwd, "./some/path")
	if src.Local.Path != expected {
		t.Errorf("expected %q, got %q", expected, src.Local.Path)
	}
}

func TestResolvePath_GitHubShorthand(t *testing.T) {
	src, err := ResolvePath("user/repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.GitHub == nil {
		t.Fatalf("expected GitHub source, got %+v", src)
	}
	if src.GitHub.Owner != "user" {
		t.Errorf("expected owner 'user', got %q", src.GitHub.Owner)
	}
	if src.GitHub.Repo != "repo" {
		t.Errorf("expected repo 'repo', got %q", src.GitHub.Repo)
	}
}

func TestResolvePath_GitHubURL(t *testing.T) {
	src, err := ResolvePath("https://github.com/user/repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.GitHub == nil {
		t.Fatalf("expected GitHub source, got %+v", src)
	}
	if src.GitHub.Owner != "user" || src.GitHub.Repo != "repo" {
		t.Errorf("expected user/repo, got %s/%s", src.GitHub.Owner, src.GitHub.Repo)
	}
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

	result = ResolveToSkillsDir(".")
	if !strings.HasSuffix(result, filepath.Join(".agents", "skills")) {
		t.Errorf("expected path ending in .agents/skills, got %q", result)
	}
}

func TestValidationError_Error(t *testing.T) {
	t.Parallel()
	e := &ValidationError{Errors: []string{"err1", "err2"}}
	msg := e.Error()
	if !strings.Contains(msg, "err1") {
		t.Errorf("expected error to contain 'err1', got %q", msg)
	}
	if !strings.Contains(msg, "err2") {
		t.Errorf("expected error to contain 'err2', got %q", msg)
	}
	if !strings.Contains(msg, "validation failed") {
		t.Errorf("expected error to contain 'validation failed', got %q", msg)
	}
}

func TestSkill_Scripts(t *testing.T) {
	t.Parallel()

	t.Run("valid map", func(t *testing.T) {
		t.Parallel()
		s := &Skill{
			Metadata: map[string]any{
				"scripts": []any{"scripts/setup.sh", "scripts/cleanup.sh"},
			},
		}
		scripts := s.Scripts()
		if len(scripts) != 2 {
			t.Fatalf("expected 2 scripts, got %d", len(scripts))
		}
		if scripts[0] != "scripts/setup.sh" {
			t.Errorf("expected 'scripts/setup.sh', got %q", scripts[0])
		}
		if scripts[1] != "scripts/cleanup.sh" {
			t.Errorf("expected 'scripts/cleanup.sh', got %q", scripts[1])
		}
	})

	t.Run("missing key", func(t *testing.T) {
		t.Parallel()
		s := &Skill{
			Metadata: map[string]any{
				"other": "value",
			},
		}
		scripts := s.Scripts()
		if len(scripts) != 0 {
			t.Errorf("expected 0 scripts, got %d", len(scripts))
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		t.Parallel()
		s := &Skill{
			Metadata: map[string]any{
				"scripts": "not-a-list",
			},
		}
		scripts := s.Scripts()
		if len(scripts) != 0 {
			t.Errorf("expected 0 scripts for wrong type, got %d", len(scripts))
		}
	})

	t.Run("nil metadata", func(t *testing.T) {
		t.Parallel()
		s := &Skill{}
		scripts := s.Scripts()
		if len(scripts) != 0 {
			t.Errorf("expected 0 scripts for nil metadata, got %d", len(scripts))
		}
	})

	t.Run("mixed types in list", func(t *testing.T) {
		t.Parallel()
		s := &Skill{
			Metadata: map[string]any{
				"scripts": []any{"scripts/valid.sh", 42, nil, "scripts/also-valid.sh"},
			},
		}
		scripts := s.Scripts()
		if len(scripts) != 2 {
			t.Fatalf("expected 2 valid scripts, got %d", len(scripts))
		}
		if scripts[0] != "scripts/valid.sh" {
			t.Errorf("expected 'scripts/valid.sh', got %q", scripts[0])
		}
		if scripts[1] != "scripts/also-valid.sh" {
			t.Errorf("expected 'scripts/also-valid.sh', got %q", scripts[1])
		}
	})
}

func TestSkill_Runtime(t *testing.T) {
	t.Parallel()

	t.Run("valid string", func(t *testing.T) {
		t.Parallel()
		s := &Skill{
			Metadata: map[string]any{
				"runtime": "python3",
			},
		}
		if s.Runtime() != "python3" {
			t.Errorf("expected 'python3', got %q", s.Runtime())
		}
	})

	t.Run("missing key", func(t *testing.T) {
		t.Parallel()
		s := &Skill{
			Metadata: map[string]any{
				"other": "value",
			},
		}
		if s.Runtime() != "" {
			t.Errorf("expected empty string, got %q", s.Runtime())
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		t.Parallel()
		s := &Skill{
			Metadata: map[string]any{
				"runtime": 42,
			},
		}
		if s.Runtime() != "" {
			t.Errorf("expected empty string for wrong type, got %q", s.Runtime())
		}
	})

	t.Run("nil metadata", func(t *testing.T) {
		t.Parallel()
		s := &Skill{}
		if s.Runtime() != "" {
			t.Errorf("expected empty string for nil metadata, got %q", s.Runtime())
		}
	})
}

func TestSkill_Dependencies(t *testing.T) {
	t.Parallel()

	t.Run("valid map", func(t *testing.T) {
		t.Parallel()
		s := &Skill{
			Metadata: map[string]any{
				"dependencies": map[string]any{
					"skills": []any{"git-commit-formatter", "pr-review"},
				},
			},
		}
		deps := s.Dependencies()
		if len(deps) != 2 {
			t.Fatalf("expected 2 deps, got %d", len(deps))
		}
		if deps[0] != "git-commit-formatter" {
			t.Errorf("expected 'git-commit-formatter', got %q", deps[0])
		}
		if deps[1] != "pr-review" {
			t.Errorf("expected 'pr-review', got %q", deps[1])
		}
	})

	t.Run("missing key", func(t *testing.T) {
		t.Parallel()
		s := &Skill{
			Metadata: map[string]any{
				"other": "value",
			},
		}
		deps := s.Dependencies()
		if deps != nil {
			t.Errorf("expected nil deps, got %v", deps)
		}
	})

	t.Run("wrong type for dependencies", func(t *testing.T) {
		t.Parallel()
		s := &Skill{
			Metadata: map[string]any{
				"dependencies": "not-a-map",
			},
		}
		deps := s.Dependencies()
		if deps != nil {
			t.Errorf("expected nil deps for wrong type, got %v", deps)
		}
	})

	t.Run("wrong type for skills", func(t *testing.T) {
		t.Parallel()
		s := &Skill{
			Metadata: map[string]any{
				"dependencies": map[string]any{
					"skills": "not-a-list",
				},
			},
		}
		deps := s.Dependencies()
		if deps != nil {
			t.Errorf("expected nil deps for wrong skills type, got %v", deps)
		}
	})

	t.Run("nil metadata", func(t *testing.T) {
		t.Parallel()
		s := &Skill{}
		deps := s.Dependencies()
		if deps != nil {
			t.Errorf("expected nil deps for nil metadata, got %v", deps)
		}
	})

	t.Run("mixed types in skills list", func(t *testing.T) {
		t.Parallel()
		s := &Skill{
			Metadata: map[string]any{
				"dependencies": map[string]any{
					"skills": []any{"valid-skill", 42, nil, "another-skill"},
				},
			},
		}
		deps := s.Dependencies()
		if len(deps) != 2 {
			t.Fatalf("expected 2 valid deps, got %d", len(deps))
		}
		if deps[0] != "valid-skill" {
			t.Errorf("expected 'valid-skill', got %q", deps[0])
		}
		if deps[1] != "another-skill" {
			t.Errorf("expected 'another-skill', got %q", deps[1])
		}
	})
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
	s := &Skill{
		Name:        "",
		Description: "A description.",
	}
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
	s := &Skill{
		Name:        "test-skill",
		Description: "",
	}
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
	s := &Skill{
		Name:        "a",
		Description: "A test.",
	}
	if err := s.Validate(); err != nil {
		t.Errorf("single char name should be valid: %v", err)
	}
}

func TestSkill_LoadFromPath_MissingFile(t *testing.T) {
	t.Parallel()
	s := &Skill{}
	err := s.LoadFromPath("/nonexistent/SKILL.md")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSkill_LoadFromPath_InvalidYAML(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	skillPath := filepath.Join(dir, SkillFileName)
	os.WriteFile(skillPath, []byte("---\ninvalid: [yaml\n---\nbody"), 0644)

	s := &Skill{}
	err := s.LoadFromPath(skillPath)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestSkill_LoadFromPath_FallbackToDirName(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "fallback-skill")
	os.MkdirAll(skillDir, 0755)
	os.WriteFile(filepath.Join(skillDir, SkillFileName), []byte("---\ndescription: desc\n---\nbody"), 0644)

	s := &Skill{}
	if err := s.LoadFromPath(filepath.Join(skillDir, SkillFileName)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != "fallback-skill" {
		t.Errorf("expected name 'fallback-skill' from dir, got %q", s.Name)
	}
}

func TestSkill_LoadFromPath_AllFields(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	content := `---
name: full-skill
description: A full skill.
tags:
  - ai
  - cli
category: productivity
author: test-author
version: 1.0.0
metadata:
  runtime: python3
  scripts:
    - scripts/run.sh
---

# Full Skill
Body text.
`
	skillPath := filepath.Join(dir, SkillFileName)
	os.WriteFile(skillPath, []byte(content), 0644)

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
	if len(s.Tags) != 2 || s.Tags[0] != "ai" || s.Tags[1] != "cli" {
		t.Errorf("expected tags [ai, cli], got %v", s.Tags)
	}
	if s.Category != "productivity" {
		t.Errorf("expected category 'productivity', got %q", s.Category)
	}
	if s.Author != "test-author" {
		t.Errorf("expected author 'test-author', got %q", s.Author)
	}
	if s.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %q", s.Version)
	}
	if s.Path != skillPath {
		t.Errorf("expected path %q, got %q", skillPath, s.Path)
	}
	if s.Content != "# Full Skill\nBody text." {
		t.Errorf("expected content '# Full Skill\nBody text.', got %q", s.Content)
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
