package skills

import (
"os"
"os/exec"
"path/filepath"
"testing"
)

func TestInstaller_sourceDir(t *testing.T) {
t.Parallel()

t.Run("explicit source dir", func(t *testing.T) {
t.Parallel()
i := &Installer{SourceDir: "/explicit/path"}
if i.sourceDir() != "/explicit/path" {
t.Errorf("expected '/explicit/path', got %q", i.sourceDir())
}
})

t.Run("empty source dir defaults to dot", func(t *testing.T) {
t.Parallel()
i := &Installer{}
if i.sourceDir() != "." {
t.Errorf("expected '.', got %q", i.sourceDir())
}
})
}

func TestInstaller_Install(t *testing.T) {
t.Parallel()

t.Run("copies SKILL.md and scripts", func(t *testing.T) {
t.Parallel()
src := t.TempDir()
skillDir := setupSourceSkill(t, src, "my-skill", "my-skill", "A test skill.", nil)

// Skill with script metadata pointing to a file alongside SKILL.md
scriptDir := filepath.Join(skillDir, "scripts")
os.MkdirAll(scriptDir, 0755)
os.WriteFile(filepath.Join(scriptDir, "setup.sh"), []byte("#!/bin/bash\necho hello\n"), 0755)

skillContent := `---
name: my-skill
description: A test skill.
metadata:
  scripts:
    - scripts/setup.sh
---

# My Skill
`
os.WriteFile(filepath.Join(skillDir, SkillFileName), []byte(skillContent), 0644)

target := t.TempDir()
parentDir := filepath.Join(target, DefaultSkillsDir)

installer := &Installer{SourceDir: src}
if err := installer.Install("my-skill", parentDir); err != nil {
t.Fatalf("Install failed: %v", err)
}

if _, err := os.Stat(filepath.Join(parentDir, SkillsSubDir, "my-skill", SkillFileName)); os.IsNotExist(err) {
t.Fatal("SKILL.md was not copied to target")
}
if _, err := os.Stat(filepath.Join(parentDir, "scripts", "setup.sh")); os.IsNotExist(err) {
t.Fatal("script was not copied to target")
}
})

t.Run("handles missing skill", func(t *testing.T) {
t.Parallel()
src := t.TempDir()
os.MkdirAll(filepath.Join(src, DefaultSkillsDir, SkillsSubDir), 0755)

installer := &Installer{SourceDir: src}
if err := installer.Install("nonexistent", filepath.Join(t.TempDir(), DefaultSkillsDir)); err == nil {
t.Fatal("expected error for missing skill")
}
})

t.Run("handles dependency resolution", func(t *testing.T) {
t.Parallel()
src := t.TempDir()
setupSourceSkill(t, src, "dep-skill", "dep-skill", "A dependency skill.", nil)

mainContent := `---
name: main-skill
description: A main skill with a dependency.
metadata:
  dependencies:
    skills:
      - dep-skill
---

# Main Skill
`
mainDir := setupSourceSkillDir(t, src, "main-skill")
os.WriteFile(filepath.Join(mainDir, SkillFileName), []byte(mainContent), 0644)

target := t.TempDir()
parentDir := filepath.Join(target, DefaultSkillsDir)

installer := &Installer{SourceDir: src}
if err := installer.Install("main-skill", parentDir); err != nil {
t.Fatalf("Install failed: %v", err)
}

for _, name := range []string{"main-skill", "dep-skill"} {
if _, err := os.Stat(filepath.Join(parentDir, SkillsSubDir, name, SkillFileName)); os.IsNotExist(err) {
t.Fatalf("skill %q was not copied", name)
}
}
})

t.Run("deduplicates dependencies", func(t *testing.T) {
t.Parallel()
src := t.TempDir()
setupSourceSkill(t, src, "shared-dep", "shared-dep", "Shared dependency.", nil)

for _, name := range []string{"skill-a", "skill-b"} {
content := `---
name: ` + name + `
description: ` + name + `.
metadata:
  dependencies:
    skills:
      - shared-dep
---
`
d := setupSourceSkillDir(t, src, name)
os.WriteFile(filepath.Join(d, SkillFileName), []byte(content), 0644)
}

target := t.TempDir()
parentDir := filepath.Join(target, DefaultSkillsDir)

installer := &Installer{SourceDir: src}
if err := installer.Install("skill-b", parentDir); err != nil {
t.Fatalf("Install failed: %v", err)
}

for _, name := range []string{"skill-b", "shared-dep"} {
if _, err := os.Stat(filepath.Join(parentDir, SkillsSubDir, name, SkillFileName)); os.IsNotExist(err) {
t.Fatalf("skill %q was not copied", name)
}
}
})
}

func TestInstaller_installAllScripts(t *testing.T) {
t.Parallel()

t.Run("copies scripts directory", func(t *testing.T) {
t.Parallel()
src := t.TempDir()
scriptsSrc := filepath.Join(src, DefaultSkillsDir, "scripts")
os.MkdirAll(scriptsSrc, 0755)
os.WriteFile(filepath.Join(scriptsSrc, "common.sh"), []byte("#!/bin/bash\n"), 0755)

target := t.TempDir()
parentDir := filepath.Join(target, DefaultSkillsDir)

installer := &Installer{SourceDir: src}
if err := installer.installAllScripts(parentDir); err != nil {
t.Fatalf("installAllScripts failed: %v", err)
}
if _, err := os.Stat(filepath.Join(parentDir, "scripts", "common.sh")); os.IsNotExist(err) {
t.Fatal("common.sh was not copied")
}
})

t.Run("missing scripts dir is not an error", func(t *testing.T) {
t.Parallel()
installer := &Installer{SourceDir: t.TempDir()}
if err := installer.installAllScripts(filepath.Join(t.TempDir(), DefaultSkillsDir)); err != nil {
t.Fatalf("expected no error for missing scripts dir, got: %v", err)
}
})
}

func TestInstaller_installScripts(t *testing.T) {
t.Parallel()

t.Run("dedup logic", func(t *testing.T) {
t.Parallel()
src := t.TempDir()
skillDir := filepath.Join(src, DefaultSkillsDir, SkillsSubDir, "skill-a")
os.MkdirAll(filepath.Join(skillDir, "scripts"), 0755)
os.WriteFile(filepath.Join(skillDir, "scripts", "shared.sh"), []byte("#!/bin/bash\n"), 0755)

skill := &Skill{
Name:        "skill-a",
Description: "Skill A",
Metadata:    map[string]any{"scripts": []any{"scripts/shared.sh"}},
}
parentDir := filepath.Join(t.TempDir(), DefaultSkillsDir)
installer := &Installer{SourceDir: src}

if err := installer.installScripts(skill, skillDir, parentDir); err != nil {
t.Fatalf("first installScripts failed: %v", err)
}
if err := installer.installScripts(skill, skillDir, parentDir); err != nil {
t.Fatalf("second installScripts failed: %v", err)
}
if _, err := os.Stat(filepath.Join(parentDir, "scripts", "shared.sh")); os.IsNotExist(err) {
t.Fatal("script was not copied")
}
})

t.Run("missing script file", func(t *testing.T) {
t.Parallel()
skillDir := filepath.Join(t.TempDir(), "skill-x")
os.MkdirAll(skillDir, 0755)

skill := &Skill{
Name:        "skill-x",
Description: "Skill X",
Metadata:    map[string]any{"scripts": []any{"scripts/missing.sh"}},
}
installer := &Installer{SourceDir: t.TempDir()}
if err := installer.installScripts(skill, skillDir, filepath.Join(t.TempDir(), DefaultSkillsDir)); err != nil {
t.Fatalf("installScripts should not error on missing script, got: %v", err)
}
})

t.Run("no scripts in metadata", func(t *testing.T) {
t.Parallel()
target := t.TempDir()
parentDir := filepath.Join(target, DefaultSkillsDir)
installer := &Installer{SourceDir: t.TempDir()}
skill := &Skill{Name: "skill-y", Description: "Skill Y"}

if err := installer.installScripts(skill, t.TempDir(), parentDir); err != nil {
t.Fatalf("unexpected error: %v", err)
}
if _, err := os.Stat(filepath.Join(parentDir, "scripts")); !os.IsNotExist(err) {
t.Fatal("scripts dir should not be created when no scripts declared")
}
})
}

func TestCopyFile(t *testing.T) {
t.Parallel()

t.Run("normal copy", func(t *testing.T) {
t.Parallel()
src := t.TempDir()
dst := t.TempDir()
srcFile := filepath.Join(src, "test.txt")
os.WriteFile(srcFile, []byte("hello world"), 0644)

dstFile := filepath.Join(dst, "test.txt")
if err := copyFile(srcFile, dstFile); err != nil {
t.Fatalf("copyFile failed: %v", err)
}
got, _ := os.ReadFile(dstFile)
if string(got) != "hello world" {
t.Errorf("expected 'hello world', got %q", got)
}
})

t.Run("missing source", func(t *testing.T) {
t.Parallel()
if err := copyFile("/nonexistent/file.txt", filepath.Join(t.TempDir(), "out.txt")); err == nil {
t.Fatal("expected error for missing source")
}
})
}

func TestCopyDir(t *testing.T) {
t.Parallel()

t.Run("nested dirs", func(t *testing.T) {
t.Parallel()
src := t.TempDir()
dst := t.TempDir()
nestedDir := filepath.Join(src, "a", "b", "c")
os.MkdirAll(nestedDir, 0755)
os.WriteFile(filepath.Join(src, "root.txt"), []byte("root"), 0644)
os.WriteFile(filepath.Join(nestedDir, "deep.txt"), []byte("deep"), 0644)

if err := copyDir(src, dst); err != nil {
t.Fatalf("copyDir failed: %v", err)
}
rootContent, _ := os.ReadFile(filepath.Join(dst, "root.txt"))
if string(rootContent) != "root" {
t.Errorf("expected 'root', got %q", rootContent)
}
deepContent, _ := os.ReadFile(filepath.Join(dst, "a", "b", "c", "deep.txt"))
if string(deepContent) != "deep" {
t.Errorf("expected 'deep', got %q", deepContent)
}
})

t.Run("missing source", func(t *testing.T) {
t.Parallel()
if err := copyDir("/nonexistent/dir", t.TempDir()); err == nil {
t.Fatal("expected error for missing source dir")
}
})
}

// setupSourceSkill creates a skill in the source directory.
func setupSourceSkill(t *testing.T, src, dirName, name, description string, scripts []string) string {
t.Helper()
skillDir := setupSourceSkillDir(t, src, dirName)

content := "---\nname: " + name + "\ndescription: " + description + "\n"
if len(scripts) > 0 {
content += "metadata:\n  scripts:\n"
for _, s := range scripts {
content += "    - " + s + "\n"
}
}
content += "---\n\n# " + name + "\n"
os.WriteFile(filepath.Join(skillDir, SkillFileName), []byte(content), 0644)
return skillDir
}

// setupSourceSkillDir creates the directory structure for a source skill.
func setupSourceSkillDir(t *testing.T, src, dirName string) string {
t.Helper()
skillDir := filepath.Join(src, DefaultSkillsDir, SkillsSubDir, dirName)
if err := os.MkdirAll(skillDir, 0755); err != nil {
t.Fatalf("creating skill dir: %v", err)
}
return skillDir
}

func runGit(t *testing.T, dir string, args ...string) {
t.Helper()
cmd := exec.Command("git", args...)
cmd.Dir = dir
if err := cmd.Run(); err != nil {
t.Fatalf("git %v failed: %v", args, err)
}
}
