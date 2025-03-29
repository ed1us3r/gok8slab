package course

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gok8slab/internal/utils"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const coursePath string = "/tmp/gok8slab/courses/"

// Course struct
type Course struct {
	Name        string   `yaml:"name"`
	Namespace   string   `yaml:"namespace"`
	Description string   `yaml:"description"`
	Guidelines  string   `yaml:"guidelines"`
	Hint        string   `yaml:"hint"`
	Manifests   []string `yaml:"manifests"`
	FlagHash    string   `yaml:"flag_hash"`

	Status bool
	Solved bool
}

// SetCourse marks a course as currently loaded.
func SetCourse(coursePath string) (*Course, error) {
	// Extract the course directory name (e.g., "000-debug" from "path/to/000-debug/")
	courseName := filepath.Base(coursePath)

	// Preferred YAML file: <coursePath>/<courseName>.yaml
	preferredFile := filepath.Join(coursePath, courseName+".yaml")

	// Fallback YAML file: <coursePath>/course.yaml
	fallbackFile := filepath.Join(coursePath, "course.yaml")

	var courseFile string
	if _, err := os.Stat(preferredFile); err == nil {
		courseFile = preferredFile
	} else if _, err := os.Stat(fallbackFile); err == nil {
		courseFile = fallbackFile
	} else {
		return nil, fmt.Errorf("no valid course file found in %s", coursePath)
	}

	// Read the determined course YAML file
	data, err := ioutil.ReadFile(courseFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read course file: %w", err)
	}

	// Parse YAML into Course struct
	var course Course
	if err := yaml.Unmarshal(data, &course); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// Create a course.lock file in the coursePath directory
	lockFilePath := filepath.Join(coursePath, "course.lock")
	if err := ioutil.WriteFile(lockFilePath, []byte(courseFile), 0644); err != nil {
		return nil, fmt.Errorf("failed to create lock file: %w", err)
	}

	logrus.Debugf("Course %s is now set as currently loaded (file: %s)", courseName, courseFile)
	return &course, nil
}

// LoadCourse reads a YAML file and parses the course
func LoadCourse(coursePath string) (*Course, error) {
	coursePath = fmt.Sprintf("%s.yaml", coursePath)
	data, err := os.ReadFile(coursePath)
	if err != nil {
		return nil, err
	}

	var course Course
	err = yaml.Unmarshal(data, &course)
	if err != nil {
		return nil, err
	}

	return &course, nil
}

// ValidateFlag checks if the user-provided flag is correct
func ValidateFlag(coursePath, userFlag string) (bool, error) {
	course, err := LoadCourse(coursePath)
	if err != nil {
		return false, err
	}

	hasher := md5.New()
	hasher.Write([]byte(userFlag))
	userFlagHash := hex.EncodeToString(hasher.Sum(nil))

	return userFlagHash == course.FlagHash, nil
}

// GetCourse retrieves a course from the given directory path if it contains a course.yaml file.
func GetCourse(dirPath string) ([]string, error) {
	var courses []string

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			courseFile := filepath.Join(dirPath, file.Name(), "course.yaml")
			if _, err := os.Stat(courseFile); err == nil {
				courses = append(courses, file.Name())
			}
		}
	}

	return courses, nil
}

// GetCurrentCourse retrieves the currently loaded course based on the course.lock file.
func GetCurrentCourse() (*Course, error) {
	fileName := "course.lock"
	dirPath := coursePath
	var lockFilePath string

	// Walk through the directory to find course.lock
	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == fileName {
			lockFilePath = path
			return fmt.Errorf("found") // Custom error to break out
		}
		return nil
	})

	// If lockFilePath is empty, course.lock was not found
	if lockFilePath == "" {
		logrus.Debug("No course.lock file found. It seems like there was never a course started yet.")
		return nil, fmt.Errorf("could not find %s in directory %s", fileName, dirPath)
	}

	// Determine the course directory
	courseDir := filepath.Dir(lockFilePath)
	courseName := filepath.Base(courseDir) // Extract directory name

	// Construct YAML file path (e.g., "000-debug/000-debug.yaml")
	yamlFilePath := filepath.Join(courseDir, courseName+".yaml")

	// Read the YAML file
	data, err := os.ReadFile(yamlFilePath)
	if err != nil {
		logrus.Debugf("Failed to read course YAML file at: %s", yamlFilePath)
		return nil, fmt.Errorf("failed to read course YAML file: %w", err)
	}

	// Parse YAML into Course struct
	var course Course
	err = yaml.Unmarshal(data, &course)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &course, nil
}

// ListCourses scans the `/courses` directory and returns only course directories.
func ListCourses() ([]string, error) {
	var courses []string
	dirPath := coursePath
	logrus.Debug("Walking through Folder:", dirPath)
	// Read the contents of the directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	// Loop through each entry in the directory
	for _, file := range files {
		// Check if the entry is a directory
		if file.IsDir() {
			logrus.Debug(">> Found Subfolder: ", file.Name())
			courses = append(courses, file.Name())
		}
	}
	return courses, nil
}

// printStatus prints out the current lab status together with the current course Informations
func PrintStatus(course *Course, error error) {
	logrus.Debug("Generating Status badge")
	statusText := "❌ No Status"
	if error == nil {

		if course.Status {
			statusText = "✅ COMPLETED"
		}

		solvedText := "❌ Not Solved yet"
		if course.Solved {
			solvedText = "✅ Solved"
		}

		border := "##############################################"
		fmt.Println(border)
		fmt.Printf("#  📚 Course:     %s\n", course.Name)
		fmt.Printf("#  🏷  Namespace:  %s\n", course.Namespace)
		fmt.Printf("#  📝 Description: %s\n", course.Description)
		fmt.Printf("#  📜 Guidelines:  %s\n", course.Guidelines)
		fmt.Printf("#  💡 Hint:        %s\n", course.Hint)
		fmt.Printf("#  📂 Manifests:   %v\n", course.Manifests)
		fmt.Printf("#  🔑 Flag Hash:   %s\n", course.FlagHash)
		fmt.Printf("#  🚀 Status:      %s\n", statusText)
		fmt.Printf("#  🎯 Solved:      %s\n", solvedText)
		fmt.Println(border)
	} else {
		border := "##############################################"

		utils.Error(border)
		utils.Error("Currently no Course Started or something else went Wrong")
		utils.Error(border)
	}
}
