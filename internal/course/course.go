package course

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"

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

// ListCourses scans the `/courses` directory and returns only course directories.
func ListCourses(dirPath string) ([]string, error) {
	var courses []string

	// Read the contents of the directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	// Loop through each entry in the directory
	for _, file := range files {
		// Check if the entry is a directory
		if file.IsDir() {
			courses = append(courses, file.Name())
		}
	}
	return courses, nil
}
