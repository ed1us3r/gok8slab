package cmd

import (
	"fmt"
	"gok8slab/internal/course"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available courses",
	Run: func(cmd *cobra.Command, args []string) {
		courses, err := course.ListCourses()
		if err != nil {

			fmt.Println("Error listing courses:", err)
			fmt.Println("Try: gok8slab pull")
			return
		}

		fmt.Println("Available Courses:")
		for _, c := range courses {
			fmt.Println(" ->", c)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
