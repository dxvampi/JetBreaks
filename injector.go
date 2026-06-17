package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func processDirectDirectory(basePath string, magicLine string) {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != "Toolbox" && !strings.Contains(strings.ToLower(file.Name()), "backup") {
			idePath := filepath.Join(basePath, file.Name())
			fmt.Printf("[*] Found IDE directory: %s\n", file.Name())
			injectVMOptions(idePath, magicLine)
		}
	}
}

func processSnapDirectory(snapPath string, magicLine string) {
	ides, err := os.ReadDir(snapPath)
	if err != nil {
		return
	}

	for _, ideDir := range ides {
		if ideDir.IsDir() {
			internalPath := filepath.Join(snapPath, ideDir.Name(), "current", ".config", "JetBrains")
			if _, err := os.Stat(internalPath); err == nil {
				subFiles, _ := os.ReadDir(internalPath)
				for _, subFile := range subFiles {
					if subFile.IsDir() && !strings.Contains(strings.ToLower(subFile.Name()), "backup") {
						idePath := filepath.Join(internalPath, subFile.Name())
						fmt.Printf("[*] Found Snap IDE directory: %s\n", subFile.Name())
						injectVMOptions(idePath, magicLine)
					}
				}
			}
		}
	}
}

func injectVMOptions(idePath string, magicLine string) {
	files, err := os.ReadDir(idePath)
	if err != nil {
		fmt.Printf("  [-] Error reading IDE directory: %v\n", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".vmoptions") {
			filePath := filepath.Join(idePath, file.Name())
			fmt.Printf("  [>] Target found: %s\n", file.Name())

			if isAlreadyInjected(filePath, "ja-netfilter") {
				fmt.Println("  [!] JetBreaks already applied here. Skipping.")
				continue
			}

			err := appendLineToFile(filePath, magicLine)
			if err != nil {
				fmt.Printf("  [-] Failed to inject: %v\n", err)
			} else {
				fmt.Println("  [+] Successfully injected!")
			}
		}
	}
}

func isAlreadyInjected(filePath string, keyword string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), keyword) {
			return true
		}
	}
	return false
}

func appendLineToFile(filePath string, line string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\n" + line + "\n")
	return err
}
