package prompt

import (
	"fmt"
	"os"
)

func ModifyPrompt(namespace string) {
	newPrompt := fmt.Sprintf(`export PS1="(Lab: %s) $PS1"`, namespace)
	_ = os.WriteFile("/etc/profile.d/gok8slab.sh", []byte(newPrompt), 0644)
}

func ResetPrompt() {
	_ = os.Remove("/etc/profile.d/gok8slab.sh")
}

