package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	dirs := []string{
		"../../collector",
		"../../extension/alloyengine",
		"../../",
	}

	for _, dir := range dirs {
		abs, err := filepath.Abs(dir)
		if err != nil {
			log.Fatalf("failed to resolve path %s: %v", dir, err)
		}
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = abs
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		log.Printf("Running go mod tidy in %s\n", abs)
		if err := cmd.Run(); err != nil {
			log.Fatalf("go mod tidy failed in %s: %v", abs, err)
		}
	}
}
