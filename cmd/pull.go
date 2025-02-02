package cmd

import (
	"gok8slab/internal/git"
	"gok8slab/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pullCmd = &cobra.Command{
	Use:   "pull [GitHub URL]",
	Short: "Pull the latest courses from a GitHub repository",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var repoURL string
		if len(args) == 1 {
			repoURL = args[0]
		} else {
			repoURL = ""
		}

		utils.Info("Pulling courses from: " + repoURL)

		// Check if dry-run mode is enabled
		if viper.GetBool("dry-run") {
			utils.Warning("Dry run mode enabled. No changes will be made.")
			return
		}

		err := git.PullCourses(repoURL, "courses")
		if err != nil {
			utils.Error("Error pulling courses: " + err.Error())
			return
		}

		utils.Success("Courses updated successfully!")
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

