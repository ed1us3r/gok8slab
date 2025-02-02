package main

import (
	"gok8slab/cmd"
	"gok8slab/internal/config"
	"os"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Execute CLI
	if err := cmd.Execute(); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
}

