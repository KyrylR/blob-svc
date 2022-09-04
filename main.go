package main

import (
	"blob-svc/internal/cli"
	"os"
)

func main() {
	// Comment out the following line to run in the IDE
	// os.Setenv("KV_VIPER_FILE", "config.yaml")
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
