package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const defaultRepo = "https://github.com/ed1us3r/gok8slab.git"
const repoDir = "gok8slab-repo"
const courseSubDir = "courses"

// PullCourses syncs courses from a GitHub repository
func PullCourses(repoURL, targetDir string) error {
	if repoURL == "" {
		repoURL = defaultRepo // Use default if none provided
	}

	// Check if repo already exists
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		fmt.Println("📂 Cloning repository for the first time...")
		err := cloneRepo(repoURL, repoDir)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("🔄 Updating existing repository...")
		err := pullRepo(repoDir)
		if err != nil {
			return err
		}
	}

	// Copy latest courses from the cloned repo
	err := copyCourses(filepath.Join(repoDir, courseSubDir), targetDir)
	if err != nil {
		return err
	}

	fmt.Println("✅ Courses updated successfully.")
	return nil
}

// Clone full repository
func cloneRepo(repoURL, targetDir string) error {
	cmd := exec.Command("git", "clone", repoURL, targetDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Clone failed:", string(output))
		return err
	}
	fmt.Println("✅ Repository cloned successfully.")
	return nil
}

// Pull latest changes from the repository
func pullRepo(targetDir string) error {
	cmd := exec.Command("git", "-C", targetDir, "pull")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Pull failed:", string(output))
		return err
	}
	fmt.Println("✅ Repository updated.")
	return nil
}

// Copy courses from cloned repo to target directory
func copyCourses(srcDir, targetDir string) error {
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return fmt.Errorf("❌ Course directory not found in repo")
	}

	// Ensure target directory exists
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Copy files using system `cp` (faster than Go's file copying)
	cmd := exec.Command("cp", "-r", srcDir+"/.", targetDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Failed to copy courses:", string(output))
		return err
	}

	fmt.Println("✅ Courses copied to:", targetDir)
	return nil
}

