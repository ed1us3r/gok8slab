package k8s

import (
	"fmt"
	"os/exec"

	"gok8slab/internal/course"
	"gok8slab/internal/utils"
	"github.com/spf13/viper"
)

func DeployCourse(c *course.Course, isOpenShift bool) error {
	if viper.GetBool("dry-run") {
		utils.Warning("Dry run mode enabled. No resources will be created.")
		fmt.Println("Would apply the following manifests:")
		for _, manifest := range c.Manifests {
			fmt.Println(" - " + manifest)
		}
		return nil
	}

	for _, manifest := range c.Manifests {
		manifestPath := fmt.Sprintf("courses/%s", manifest)
		var cmd *exec.Cmd

		if isOpenShift {
			utils.Info("Using OpenShift client (oc apply)")
			cmd = exec.Command("oc", "apply", "-f", manifestPath, "-n", c.Namespace)
		} else {
			utils.Info("Using Kubernetes client (kubectl apply)")
			cmd = exec.Command("kubectl", "apply", "-f", manifestPath, "-n", c.Namespace)
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			utils.Error("Deployment failed: " + err.Error())
			return err
		}
		utils.Success("Deployment successful: " + manifest)
		fmt.Printf("%s\n", output)
	}

	return nil
}

