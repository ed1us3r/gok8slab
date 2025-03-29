package cmd

import (
	"fmt"
	"gok8slab/internal/course"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Retrieve Status Information about the lab Status.",
	Run: func(cmd *cobra.Command, args []string) {
		courses, err := course.ListCourses()
		if err != nil {

			fmt.Println("Error Retrieving status of Lab:", err)
			fmt.Println("Try: gok8slab help")
			return
		}

		currentCourse, err := course.GetCurrentCourse()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		course.PrintStatus(currentCourse, err)
		fmt.Println("Other Available Courses:")
		for _, c := range courses {
			fmt.Println(" -", c)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
