package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLoader(t *testing.T) {
	l := NewLoader("/some/path")
	if l.RootPath != "/some/path" {
		t.Errorf("expected RootPath '/some/path', got %q", l.RootPath)
	}
}

func TestLoader_LoadAll(t *testing.T) {
	t.Parallel()

	t.Run("multi-skill dir", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		skillsDir := filepath.Join(root, DefaultSkillsDir, SkillsSubDir)

		// Create two valid skills
		createSkillDir(t, skillsDir, "alpha", "alpha", "An alpha skill.")
		createSkillDir(t, skillsDir, "beta", "beta", "A beta skill.")

		l := NewLoader(root)
		skills, err := l.LoadAll()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(skills) != 2 {
			t.Fatalf("expected 2 skills, got %d", len(skills))
		}
		names := make(map[string]bool)
		for _, s := range skills {
			names[s.Name] = true
		}
		if !names["alpha"] {
			t.Error("expected skill 'alpha'")
		}
		if !names["beta"] {
			t.Error("expected skill 'beta'")
		}
	})

	t.Run("valid and invalid entries", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		skillsDir := filepath.Join(root, DefaultSkillsDir, SkillsSubDir)

		// Valid skill
		createSkillDir(t, skillsDir, "valid", "valid", "A valid skill.")

		// Directory without SKILL.md — should be skipped silently
		os.MkdirAll(filepath.Join(skillsDir, "no-file"), 0755)

		// Directory with invalid YAML — should warn and skip
		invalidDir := filepath.Join(skillsDir, "invalid")
		os.MkdirAll(invalidDir, 0755)
		os.WriteFile(filepath.Join(invalidDir, SkillFileName), []byte("---\nbad yaml: [\n---\nbody"), 0644)

		l := NewLoader(root)
		skills, err := l.LoadAll()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(skills) != 1 {
			t.Fatalf("expected 1 skill, got %d", len(skills))
		}
		if skills[0].Name != "valid" {
			t.Errorf("expected skill 'valid', got %q", skills[0].Name)
		}
	})

	t.Run("empty dir", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		skillsDir := filepath.Join(root, DefaultSkillsDir, SkillsSubDir)
		os.MkdirAll(skillsDir, 0755)

		l := NewLoader(root)
		skills, err := l.LoadAll()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(skills) != 0 {
			t.Fatalf("expected 0 skills, got %d", len(skills))
		}
	})

	t.Run("missing dir", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		// Do NOT create .agents/skills

		l := NewLoader(root)
		_, err := l.LoadAll()
		if err == nil {
			t.Fatal("expected error for missing skills directory")
		}
	})

	t.Run("non-directory entries skipped", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		skillsDir := filepath.Join(root, DefaultSkillsDir, SkillsSubDir)
		os.MkdirAll(skillsDir, 0755)

		// Create a file (not directory) in skills dir
		os.WriteFile(filepath.Join(skillsDir, "not-a-dir.md"), []byte("content"), 0644)

		l := NewLoader(root)
		skills, err := l.LoadAll()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(skills) != 0 {
			t.Fatalf("expected 0 skills, got %d", len(skills))
		}
	})
}

// createSkillDir is a helper that creates a skill directory with a SKILL.md file.
func createSkillDir(t *testing.T, parentDir, name, skillName, description string) {
	t.Helper()
	skillDir := filepath.Join(parentDir, name)
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("creating skill dir: %v", err)
	}
	content := "---\nname: " + skillName + "\ndescription: " + description + "\n---\n\n# " + name + "\n"
	if err := os.WriteFile(filepath.Join(skillDir, SkillFileName), []byte(content), 0644); err != nil {
		t.Fatalf("writing SKILL.md: %v", err)
	}
}
