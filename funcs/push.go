package funcs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Push(filename string) {
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		fmt.Println("Initialize the repo using: codat init")
		return
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	configPath := filepath.Join(homeDir, configDir)
	configFilePath := filepath.Join(configPath, configFile)

	token, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Error reading token:", err)
		return
	}
	tokenStr := strings.TrimSpace(strings.TrimPrefix(string(token), tokenKey))
	if tokenStr == "" {
		fmt.Println("Token is empty. Set the token using: codat set-token <token>")
		return
	}

	originPath := filepath.Join(repoDir, "HEAD")
	origin, err := os.ReadFile(originPath)
	if err != nil {
		fmt.Println("Error reading the HEAD file:", err)
		return
	}
	originStr := strings.TrimSpace(string(origin))
	if originStr == "" {
		fmt.Println("Origin is empty. Set the origin using: codat set-origin <codatLink>")
		return
	}

	parts := strings.Split(originStr, "/")
	codatID := parts[len(parts)-1]

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", filename)
		return
	}

	fmt.Println("Pushing file:", filename)
	fmt.Println("Using token:", tokenStr)
	fmt.Println("Pushing to:", originStr)
	fmt.Println("File content:\n", string(content))

	apiURL := fmt.Sprintf("http://localhost:3000/api/profile/verify-token?codatId=%s", codatID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating API request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+tokenStr)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making API request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading API response:", err)
		return
	}

	var apiResponse map[string]interface{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Println("Error parsing API response:", err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println("API Error:", apiResponse["error"])
		return
	}

	fmt.Println("API Verification Successful:", apiResponse["message"])
}
