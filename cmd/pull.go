package cmd

import (
	"fmt"

	"gok8slab/internal/git"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [GitHub URL]",
	Short: "Pull the latest courses from a GitHub repository",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var repoURL string

		if len(args) == 1 {
			repoURL = args[0] // Use provided GitHub repo
		} else {
			repoURL = "" // Use default repo (gok8slab)
		}

		fmt.Println("Pulling courses from:", repoURL)

		err := git.PullCourses(repoURL, "courses")
		if err != nil {
			fmt.Println("❌ Error pulling courses:", err)
			return
		}

		fmt.Println("✅ Courses updated successfully!")
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

