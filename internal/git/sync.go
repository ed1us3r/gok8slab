package git

import (
	"fmt"
	"gok8slab/internal/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Constants
const (
	defaultRepo  = "https://github.com/ed1us3r/gok8slab.git"
	repoDir      = "gok8slab-repo"
	courseSubDir = "courses"
)

// PullCourses syncs courses from a GitHub repository
func PullCourses(repoURL, targetDir string) error {
	if repoURL == "" {
		repoURL = defaultRepo // Use default if none provided
	}

	if err := syncRepo(repoURL); err != nil {
		return err
	}

	// Copy the latest courses from the cloned repo
	coursesSrcDir := filepath.Join(repoDir, courseSubDir)

	if err := copyCourses(coursesSrcDir, targetDir); err != nil {
		return err
	}

	utils.Success("✅ Courses updated successfully.")
	return nil
}

// syncRepo handles cloning or pulling the repository
func syncRepo(repoURL string) error {
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		utils.Info("📂 Cloning repository for the first time...")
		return cloneRepo(repoURL, repoDir)
	} else {
		utils.Info("🔄 Updating existing repository...")
		return pullRepo(repoDir)
	}
}

// Clone full repository
func cloneRepo(repoURL, targetDir string) error {
	cmd := exec.Command("git", "clone", repoURL, targetDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.Error("❌ Clone failed: \n")
		fmt.Printf("[Error] %s ", output)

		return err
	}
	utils.Success("✅ Repository cloned successfully.")
	return nil
}

// Pull latest changes from the repository
func pullRepo(targetDir string) error {
	cmd := exec.Command("git", "-C", targetDir, "pull")
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.Error("❌ Pull failed: \n")

		fmt.Printf("[Error] %v", output)
		return err
	}
	utils.Success("✅ Repository updated.")
	return nil
}

// Copy courses from cloned repo to target directory
func copyCourses(srcDir, targetDir string) error {
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(srcDir)
	if err != nil {
		utils.Error("❌ Failed to Access the Course directory: \n")

		return fmt.Errorf("[Error] %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if err := copyCourseFiles(filepath.Join(srcDir, entry.Name()), targetDir, entry.Name()); err != nil {
				return err
			}
		}
	}
	return nil
}

// Copies course files and Kubernetes manifests associated with the course
func copyCourseFiles(courseDir, targetDir, courseName string) error {
	courseFiles, err := ioutil.ReadDir(courseDir)
	if err != nil {
		utils.Error("❌ Failed to read course folder: ")
		return fmt.Errorf("[Error] %v", err)
	}
	targetCourseDir := filepath.Join(courseDir, courseName)
	logrus.Debugf("Mkdir: %s", targetCourseDir)
	if err := os.MkdirAll(targetCourseDir, os.ModePerm); err != nil {
		return err
	}

	utils.Info("💤 Updating and copying Course Files...")
	utils.Info("💤 Currently Loading Course: " + courseName)
	for _, courseFile := range courseFiles {
		srcPath := filepath.Join(courseDir, courseFile.Name())
		dstPath := filepath.Join(targetDir, courseName, courseFile.Name())
		logrus.Debugf("Moving Files: %s -> %s", srcPath, dstPath)
		if filepath.Ext(courseFile.Name()) == ".yaml" {
			// If it's a YAML file, copy it
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		} else if courseFile.IsDir() {
			// Handle Kubernetes manifests directory
			if err := copyKubernetesManifests(srcPath, targetDir, courseName); err != nil {
				return err
			}
		}
	}

	utils.Success("Fully loaded Course: " + courseName)
	return nil
}

// Copy a single file
func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(dst, input, os.ModePerm); err != nil {
		return err
	}
	utils.Info("   ✅ Copied file:" + dst)
	return nil
}

// Copy Kubernetes manifests from a directory
func copyKubernetesManifests(srcDir, targetDir, courseName string) error {
	manifestsDir := filepath.Join(srcDir, courseName) // Assuming manifest folder mirrors course name
	if _, err := os.Stat(manifestsDir); err != nil {
		return nil // Skip if the directory doesn't exist
	}

	files, err := ioutil.ReadDir(manifestsDir)
	if err != nil {
		utils.Error("❌ Failed to read Kubernetes manifests directory")
		return fmt.Errorf("❌ Failed to read Kubernetes manifests directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" {
			err := copyFile(filepath.Join(manifestsDir, file.Name()), filepath.Join(targetDir, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
