package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Remote URLs for the transparent FOSS assets
// (Once you create your GitHub repo, replace these with your actual RAW URLs)
const (
	baseURL  = "https://raw.githubusercontent.com/your-username/jetbreaks/main/assets/"
	agentURL = baseURL + "ja-netfilter.jar"
)

// List of essential configuration files needed by ja-netfilter
var configFiles = []string{
	"dns.conf",
	"url.conf",
}

func setupEnvironment(dirPath string, agentPath string) error {
	// 1. Create main ~/.jetbreaks directory
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create main directory: %v", err)
		}
		fmt.Println("  [+] Created directory: ~/.jetbreaks")
	}

	// 2. Create the internal config directory (~/.jetbreaks/config)
	configDirPath := filepath.Join(dirPath, "config")
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
		fmt.Println("  [+] Created subdirectory: ~/.jetbreaks/config")
	}

	// 3. Handle the ja-netfilter.jar download
	if _, err := os.Stat(agentPath); os.IsNotExist(err) {
		fmt.Println("  [>] Downloading ja-netfilter.jar...")
		if err := downloadFile(agentPath, agentURL); err != nil {
			fmt.Printf("  [!] Agent download failed (%v). Creating mock placeholder.\n", err)
			_ = createPlaceholderFile(agentPath, "MOCK JAR")
		} else {
			fmt.Println("  [+] Agent downloaded successfully.")
		}
	}

	// 4. Handle the configuration files loop (.conf)
	for _, confName := range configFiles {
		localConfPath := filepath.Join(configDirPath, confName)
		remoteConfURL := baseURL + "config/" + confName

		if _, err := os.Stat(localConfPath); os.IsNotExist(err) {
			fmt.Printf("  [>] Downloading configuration: %s...\n", confName)
			if err := downloadFile(localConfPath, remoteConfURL); err != nil {
				fmt.Printf("  [!] %s download failed (%v). Creating empty placeholder.\n", confName, err)
				_ = createPlaceholderFile(localConfPath, "[dns]\nEQUAL,jetbrains.com") // Safe testing rule
			} else {
				fmt.Printf("  [+] %s downloaded successfully.\n", confName)
			}
		}
	}

	return nil
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func createPlaceholderFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}
