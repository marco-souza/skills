package skills

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Installer copies skills to target projects.
type Installer struct {
	SourceDir string
}

func (i *Installer) sourceDir() string {
	if i.SourceDir != "" {
		return i.SourceDir
	}
	return "."
}

// Install copies a single skill to the target project.
// parentDir is the directory that contains the skills/ subdirectory
// (e.g., ".agents" or ".opencode").
func (i *Installer) Install(skillName, parentDir string) error {
	skillsDir := filepath.Join(i.sourceDir(), defaultSkillsDir, skillsSubDir)
	skillSrc := filepath.Join(skillsDir, skillName)

	if _, err := os.Stat(skillSrc); os.IsNotExist(err) {
		return fmt.Errorf("skill %q not found in %s", skillName, skillsDir)
	}

	skillDest := filepath.Join(parentDir, skillsSubDir, skillName)

	if err := os.MkdirAll(filepath.Dir(skillDest), 0755); err != nil {
		return fmt.Errorf("creating target directory: %w", err)
	}

	if err := copyDir(skillSrc, skillDest); err != nil {
		return fmt.Errorf("copying skill: %w", err)
	}

	// Load the skill to resolve script dependencies, then copy them.
	skillPath := filepath.Join(skillSrc, skillFileName)
	skill := &Skill{}
	if err := skill.LoadFromPath(skillPath); err == nil {
		if err := i.installScripts(skill, skillSrc, parentDir); err != nil {
			return fmt.Errorf("installing scripts for %q: %w", skillName, err)
		}
	}

	fmt.Printf("Installed skill %q to %s\n", skillName, skillDest)
	return nil
}

// InstallFromGitHub clones a GitHub repo and installs all its skills to the target.
// parentDir is the directory that contains the skills/ subdirectory.
func (i *Installer) InstallFromGitHub(src *GitHubSource, parentDir string) error {
	repoDir, cleanup, err := CloneRepo(src)
	if err != nil {
		return err
	}
	defer cleanup()

	skillsPath := filepath.Join(repoDir, defaultSkillsDir, skillsSubDir)
	if _, err := os.Stat(skillsPath); os.IsNotExist(err) {
		skillsPath = repoDir
	}

	entries, err := os.ReadDir(skillsPath)
	if err != nil {
		return fmt.Errorf("reading skills in cloned repo: %w", err)
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		skillSrc := filepath.Join(skillsPath, entry.Name())
		skillDest := filepath.Join(parentDir, skillsSubDir, entry.Name())

		if err := os.MkdirAll(filepath.Dir(skillDest), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "warning: %v\n", err)
			continue
		}
		if err := copyDir(skillSrc, skillDest); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to copy %s: %v\n", entry.Name(), err)
			continue
		}

		// Load the skill to resolve script dependencies, then copy them.
		skillPath := filepath.Join(skillSrc, skillFileName)
		skill := &Skill{}
		if err := skill.LoadFromPath(skillPath); err == nil {
			if err := i.installScripts(skill, skillSrc, parentDir); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to install scripts for %s: %v\n", entry.Name(), err)
			}
		}

		count++
	}

	fmt.Printf("Installed %d skills from %s to %s\n", count, src.SSHURL, parentDir)
	return nil
}

// copyFile copies a single file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// installScripts copies script dependencies declared in a skill's metadata
// to the target project's scripts directory (parentDir/scripts/).
// Scripts that already exist in the destination are skipped.
func (i *Installer) installScripts(skill *Skill, skillSrcDir, parentDir string) error {
	scriptPaths := skill.Scripts()
	if len(scriptPaths) == 0 {
		return nil
	}

	scriptsDestDir := filepath.Join(parentDir, "scripts")
	if err := os.MkdirAll(scriptsDestDir, 0755); err != nil {
		return fmt.Errorf("creating scripts directory: %w", err)
	}

	for _, relPath := range scriptPaths {
		// Resolve relative to the skill's source directory
		absSrc := filepath.Clean(filepath.Join(skillSrcDir, relPath))

		if _, err := os.Stat(absSrc); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "warning: script %q not found at %s, skipping\n", relPath, absSrc)
			continue
		}

		destPath := filepath.Join(scriptsDestDir, filepath.Base(relPath))
		if _, err := os.Stat(destPath); err == nil {
			continue // already installed by another skill
		}

		if err := copyFile(absSrc, destPath); err != nil {
			return fmt.Errorf("copying script %s: %w", relPath, err)
		}
		fmt.Printf("Installed script %q to %s\n", filepath.Base(relPath), destPath)
	}

	return nil
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() { _ = srcFile.Close() }()

		dstFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer func() { _ = dstFile.Close() }()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}
