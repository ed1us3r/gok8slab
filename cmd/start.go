package cmd

import (
	"fmt"
	"gok8slab/internal/course"
	"gok8slab/internal/k8s"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:   "start [course.yaml]",
	Short: "Start a Kubernetes learning course",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		coursePath := filepath.Join("/tmp/gok8slab/courses", args[0], args[0])

		courseData, err := course.LoadCourse(coursePath)
		if err != nil {
			fmt.Println("Failed to load course:", err)
			return
		}

		fmt.Println("Starting Course:", courseData.Name)
		fmt.Println("Description:", courseData.Description)
		fmt.Println("Guidelines:", courseData.Guidelines)

		err = k8s.DeployCourse(courseData, viper.GetBool("openshift"))
		if err != nil {
			fmt.Println("Failed to start course:", err)
			return
		}
		course.SetCourse(coursePath)
		fmt.Println("Course deployed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
