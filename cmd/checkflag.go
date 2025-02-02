package cmd

import (
	"fmt"
	"path/filepath"

	"gok8slab/internal/course"
	"github.com/spf13/cobra"
)

var checkFlagCmd = &cobra.Command{
	Use:   "checkflag [course.yaml] [flag]",
	Short: "Check if the provided flag is correct",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		coursePath := filepath.Join("courses", args[0])
		userFlag := args[1]

		valid, err := course.ValidateFlag(coursePath, userFlag)
		if err != nil {
			fmt.Println("Error checking flag:", err)
			return
		}

		if valid {
			fmt.Println("✅ Correct flag! Well done! 🎉")
		} else {
			fmt.Println("❌ Incorrect flag. Keep investigating!")
		}
	},
}

func init() {
	rootCmd.AddCommand(checkFlagCmd)
}

