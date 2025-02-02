package k8s

import (
	"fmt"
	"os/exec"

	"gok8slab/internal/course"
	"github.com/sirupsen/logrus"
)

func DeployCourse(c *course.Course, isOpenShift bool) error {
	var cmd *exec.Cmd

	for _, manifest := range c.Manifests {
		manifestPath := fmt.Sprintf("courses/%s", manifest)

		if isOpenShift {
			logrus.Info("Using OpenShift client (oc apply)")
			cmd = exec.Command("oc", "apply", "-f", manifestPath, "-n", c.Namespace)
		} else {
			logrus.Info("Using Kubernetes client (kubectl apply)")
			cmd = exec.Command("kubectl", "apply", "-f", manifestPath, "-n", c.Namespace)
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			logrus.Error("Deployment failed:", err)
			return err
		}
		logrus.Debug(string(output))
	}

	return nil
}

