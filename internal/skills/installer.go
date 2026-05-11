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
	// SourceDir is the directory containing the source skills. If empty, the current directory is used.
	SourceDir string
}

func (i *Installer) sourceDir() string {
	if i.SourceDir != "" {
		return i.SourceDir
	}
	return "."
}

// Install copies a single skill and all its skill dependencies to the target project.
// The parentDir argument specifies the directory that contains the skills/ subdirectory (e.g., ".agents").
func (i *Installer) Install(skillName, parentDir string) error {
	return i.installWithTracking(skillName, parentDir, make(map[string]bool))
}

// installWithTracking installs a skill and its dependencies, skipping
// already-installed skills via the installed set.
func (i *Installer) installWithTracking(skillName, parentDir string, installed map[string]bool) error {
	if installed[skillName] {
		return nil // already installed as a dependency
	}

	skillsDir := filepath.Join(i.sourceDir(), DefaultSkillsDir, SkillsSubDir)
	skillSrc := filepath.Join(skillsDir, skillName)

	if _, err := os.Stat(skillSrc); os.IsNotExist(err) {
		return fmt.Errorf("skill %q not found in %s", skillName, skillsDir)
	}

	// Load the skill to extract metadata before copying
	skillPath := filepath.Join(skillSrc, SkillFileName)
	skill := &Skill{}
	if err := skill.LoadFromPath(skillPath); err != nil {
		return fmt.Errorf("loading skill %q: %w", skillName, err)
	}

	// Install dependency skills first (recursive)
	for _, depName := range skill.Dependencies() {
		if installed[depName] {
			continue
		}
		if err := i.installWithTracking(depName, parentDir, installed); err != nil {
			return fmt.Errorf("installing dependency %q for %q: %w", depName, skillName, err)
		}
	}

	// Copy the skill itself
	skillDest := filepath.Join(parentDir, SkillsSubDir, skillName)
	if err := os.MkdirAll(filepath.Dir(skillDest), 0755); err != nil {
		return fmt.Errorf("creating target directory: %w", err)
	}
	if err := copyDir(skillSrc, skillDest); err != nil {
		return fmt.Errorf("copying skill: %w", err)
	}

	// Install script dependencies declared by this skill
	if err := i.installScripts(skill, skillSrc, parentDir); err != nil {
		return fmt.Errorf("installing scripts for %q: %w", skillName, err)
	}

	installed[skillName] = true
	fmt.Printf("Installed skill %q to %s\n", skillName, skillDest)
	return nil
}

// InstallFromGitHub clones a GitHub repository and installs every skill it contains into the target project.
// The parentDir argument specifies the directory that contains the skills/ subdirectory.
func (i *Installer) InstallFromGitHub(src *GitHubSource, parentDir string) error {
	repoDir, cleanup, err := CloneRepo(src)
	if err != nil {
		return err
	}
	defer cleanup()

	skillsPath := filepath.Join(repoDir, DefaultSkillsDir, SkillsSubDir)
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
		skillDest := filepath.Join(parentDir, SkillsSubDir, entry.Name())

		if err := os.MkdirAll(filepath.Dir(skillDest), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "warning: %v\n", err)
			continue
		}
		if err := copyDir(skillSrc, skillDest); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to copy %s: %v\n", entry.Name(), err)
			continue
		}

		// Load the skill to resolve script dependencies, then copy them.
		skillPath := filepath.Join(skillSrc, SkillFileName)
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
func copyFile(src, dst string) (err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { if cerr := dstFile.Close(); cerr != nil && err == nil { err = cerr } }()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// InstallAll installs every skill from the source directory to the target project.
// It also copies the entire scripts/ directory so all shared scripts are available.
// The parentDir argument specifies the directory that contains the skills/ subdirectory (e.g., ".agents").
func (i *Installer) InstallAll(parentDir string) error {
	skillsPath := filepath.Join(i.sourceDir(), DefaultSkillsDir, SkillsSubDir)
	if _, err := os.Stat(skillsPath); os.IsNotExist(err) {
		// Fall back to source root (mirrors InstallFromGitHub behavior)
		skillsPath = i.sourceDir()
	}

	entries, err := os.ReadDir(skillsPath)
	if err != nil {
		return fmt.Errorf("reading skills directory %q: %w", skillsPath, err)
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if err := i.Install(entry.Name(), parentDir); err != nil {
			fmt.Fprintf(os.Stderr, "warning: %v\n", err)
			continue
		}
		count++
	}

	if count == 0 {
		return fmt.Errorf("no skills found in %s", skillsPath)
	}

	// Copy the entire scripts directory so all shared scripts are available.
	if err := i.InstallAllScripts(parentDir); err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}

	fmt.Printf("Installed %d skills to %s\n", count, filepath.Join(parentDir, "skills"))
	return nil
}

// InstallAllScripts copies the entire scripts/ directory from the source to the target project.
// It should be called after installing all skills to ensure every shared script is available.
func (i *Installer) InstallAllScripts(parentDir string) error {
	scriptsSrc := filepath.Join(i.sourceDir(), DefaultSkillsDir, "scripts")
	if _, err := os.Stat(scriptsSrc); os.IsNotExist(err) {
		return nil // no scripts directory — not an error
	}

	scriptsDest := filepath.Join(parentDir, "scripts")
	if err := os.MkdirAll(scriptsDest, 0755); err != nil {
		return fmt.Errorf("creating scripts directory: %w", err)
	}

	if err := copyDir(scriptsSrc, scriptsDest); err != nil {
		return fmt.Errorf("copying scripts: %w", err)
	}

	fmt.Printf("Installed all scripts to %s\n", scriptsDest)
	return nil
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
		defer srcFile.Close()

		dstFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer func() { if cerr := dstFile.Close(); cerr != nil && err == nil { err = cerr } }()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}
