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
	local, ok := src.(*LocalSource)
	if !ok {
		t.Fatalf("expected LocalSource, got %T", src)
	}
	if local.Path != "/absolute/path" {
		t.Errorf("expected '/absolute/path', got %q", local.Path)
	}
}

func TestResolvePath_Relative(t *testing.T) {
	src, err := ResolvePath("./some/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	local, ok := src.(*LocalSource)
	if !ok {
		t.Fatalf("expected LocalSource, got %T", src)
	}
	cwd, _ := os.Getwd()
	expected := filepath.Join(cwd, "./some/path")
	if local.Path != expected {
		t.Errorf("expected %q, got %q", expected, local.Path)
	}
}

func TestResolvePath_GitHubShorthand(t *testing.T) {
	src, err := ResolvePath("user/repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gh, ok := src.(*GitHubSource)
	if !ok {
		t.Fatalf("expected GitHubSource, got %T", src)
	}
	if gh.Owner != "user" {
		t.Errorf("expected owner 'user', got %q", gh.Owner)
	}
	if gh.Repo != "repo" {
		t.Errorf("expected repo 'repo', got %q", gh.Repo)
	}
}

func TestResolvePath_GitHubURL(t *testing.T) {
	src, err := ResolvePath("https://github.com/user/repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gh, ok := src.(*GitHubSource)
	if !ok {
		t.Fatalf("expected GitHubSource, got %T", src)
	}
	if gh.Owner != "user" || gh.Repo != "repo" {
		t.Errorf("expected user/repo, got %s/%s", gh.Owner, gh.Repo)
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
