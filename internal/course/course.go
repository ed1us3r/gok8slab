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
	courseFile := fmt.Sprintf("%s.yaml", coursePath)

	data, err := ioutil.ReadFile(courseFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read course file: %w", err)
	}

	var course Course
	if err := yaml.Unmarshal(data, &course); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// Create a course.lock file in the coursePath directory
	lockFilePath := filepath.Join(filepath.Dir(coursePath), "course.lock")
	if err := ioutil.WriteFile(lockFilePath, []byte(courseFile), 0644); err != nil {
		return nil, fmt.Errorf("failed to create lock file: %w", err)
	}

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
func GetCurrentCourse(dirPath string) (*Course, error) {
	lockFilePath := filepath.Join(dirPath, "course.lock")
	data, err := os.ReadFile(lockFilePath)
	if err != nil {
		fmt.Errorf("We could not detect a cource.lock File at the Course-Directory: %w", dirPath)
		logrus.Debug("It seems like there was never a course started yet. Make sure to run one first.")
		return nil, fmt.Errorf("failed to read course.lock file: %w", err)
	}
	var course Course
	logrus.Debug("Loaded Course:", course)
	err = yaml.Unmarshal(data, &course)
	if err != nil {
		return nil, err
	}

	return &course, nil
}

// ListCourses scans the `/courses` directory and returns only course directories.
func ListCourses(dirPath string) ([]string, error) {
	var courses []string

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
	statusText := "❌ FAILED"
	if error != nil {

		if course.Status {
			statusText = "✅ COMPLETED"
		}

		solvedText := "❌ Not Solved"
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
