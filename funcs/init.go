package funcs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	repoDir    = ".codat"
	configDir  = ".codat"
	configFile = "config"
	tokenKey   = "token="
)

func Init() {
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		if err := os.Mkdir(repoDir, 0755); err != nil {
			fmt.Println("Error creating repository directory:", err)
			return
		}
		if err := os.WriteFile(filepath.Join(repoDir, "HEAD"), []byte(""), 0644); err != nil {
			fmt.Println("Error creating HEAD file:", err)
			return
		}
		fmt.Println("Initialized empty repository in", repoDir)
	} else {
		fmt.Println("Repository already initialized.")
	}
}

func SetOrigin(codatLink string) {
	codatLink = strings.TrimSpace(codatLink)
	prefix := "http://localhost:3000/codat/"

	if !strings.HasPrefix(codatLink, prefix) {
		fmt.Println("Invalid codat link")
		return
	}
	if err := os.WriteFile(filepath.Join(repoDir, "HEAD"), []byte(codatLink), 0644); err != nil {
		fmt.Println("Error writing to HEAD file:", err)
	}
}

func SetToken(token string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	configPath := filepath.Join(homeDir, configDir)
	configFilePath := filepath.Join(configPath, configFile)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.Mkdir(configPath, 0755); err != nil {
			fmt.Println("Error creating config directory:", err)
			return
		}
	}

	content := tokenKey + token + "\n"
	if err := os.WriteFile(configFilePath, []byte(content), 0644); err != nil {
		fmt.Println("Error writing token to config file:", err)
		return
	}

	fmt.Println("Token set successfully in", configFilePath)
}
