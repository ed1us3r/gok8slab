package k8s

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

func Cleanup(isOpenShift bool) error {
	var cmd *exec.Cmd

	if isOpenShift {
		logrus.Info("Using OpenShift client (oc delete)")
		cmd = exec.Command("oc", "delete", "-f", "/tmp/k8s-labs")
	} else {
		logrus.Info("Using Kubernetes client (kubectl delete)")
		cmd = exec.Command("kubectl", "delete", "-f", "/tmp/k8s-labs")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Error("Failed to cleanup:", err)
		return err
	}

	logrus.Debug("Cleanup Output:", string(output))
	return nil
}

