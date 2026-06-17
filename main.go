package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("=== JetBreaks v0.1.0 - Starting Smart Injector ===")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("[-] Error getting Home directory: %v\n", err)
		return
	}

	// Define our local workspace directory
	jetBreaksDir := filepath.Join(homeDir, ".jetbreaks")
	agentPath := filepath.Join(jetBreaksDir, "ja-netfilter.jar")
	magicLine := fmt.Sprintf("-javaagent:%s=jetbrains", agentPath)

	// 1. Setup the local environment (Folders & Real HTTP Downloads)
	fmt.Println("[*] Setting up local environment...")
	err = setupEnvironment(jetBreaksDir, agentPath)
	if err != nil {
		fmt.Printf("[-] Environment setup failed: %v\n", err)
		return
	}

	// 2. Scan and Inject JetBrains IDEs
	possiblePaths := []string{
		filepath.Join(homeDir, ".config", "JetBrains"),
		filepath.Join(homeDir, "snap"),
		filepath.Join(homeDir, ".var", "app"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		if filepath.Base(path) == "JetBrains" {
			processDirectDirectory(path, magicLine)
		}
		if filepath.Base(path) == "snap" {
			processSnapDirectory(path, magicLine)
		}
	}
}
