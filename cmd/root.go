package cmd

import (
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
)

var rootCmd = &cobra.Command{
	Use:   "gok8slab",
	Short: "gok8slab - Kubernetes Capture The Flag Learning Environment",
	Long: `A CLI tool to set up Kubernetes CTF-style challenges.
It deploys lessons to a Kubernetes or OpenShift cluster and cleans up after completion.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set log level based on flag
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Debug("Debug mode enabled")
		} else if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			logrus.SetLevel(logrus.InfoLevel)
		} else {
			logrus.SetLevel(logrus.WarnLevel)
		}

		// Detect OpenShift if flag is not explicitly set
		if !viper.IsSet("openshift") {
			isOpenShift := detectOpenShift()
			viper.Set("openshift", isOpenShift)
		}

		logrus.Infof("Using OpenShift: %v", viper.GetBool("openshift"))
	},
}

// Execute CLI
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().Bool("openshift", false, "Force OpenShift mode")

	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("openshift", rootCmd.PersistentFlags().Lookup("openshift"))
}

// Detect OpenShift by checking if the "oc" command exists
func detectOpenShift() bool {
	_, err := exec.LookPath("oc")
	return err == nil
}

