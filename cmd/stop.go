package cmd

import (
	"fmt"
	"gok8slab/internal/k8s"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Cleanup the lab environment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stopping lab environment...")
		err := k8s.Cleanup(viper.GetBool("openshift"))
		if err != nil {
			fmt.Println("Failed to cleanup:", err)
			return
		}
		fmt.Println("Lab cleaned up successfully.")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

