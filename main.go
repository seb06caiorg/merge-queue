package main

import (
	"os"
	"os/exec"
)

// Simple wrapper that runs the actual server application.
// The main application is located in cmd/server/main.go.
func main() {
	cmd := exec.Command("go", "run", "cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
