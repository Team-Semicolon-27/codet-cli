package funcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var mapExtensions = map[string]string{
	".go":    "Go",
	".ts":    "TypeScript",
	".js":    "JavaScript",
	".py":    "Python",
	".java":  "Java",
	".cpp":   "C++",
	".c":     "C",
	".rs":    "Rust",
	".rb":    "Ruby",
	".php":   "PHP",
	".swift": "Swift",
}

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

	ext := filepath.Ext(filename)
	language, exists := mapExtensions[ext]
	if !exists {
		language = "Unknown"
	}

	fmt.Println("Pushing file:", filename)
	fmt.Println("Detected language:", language)
	fmt.Println("Using token:", tokenStr)
	fmt.Println("Pushing to:", originStr)

	client := &http.Client{}

	patchURL := fmt.Sprintf("http://localhost:3000/api/codat/edit/%s", codatID)

	requestBody, err := json.Marshal(map[string]string{
		"code":     string(content),
		"language": language,
	})
	if err != nil {
		fmt.Println("Error creating request body:", err)
		return
	}

	patchReq, err := http.NewRequest("PATCH", patchURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating PATCH request:", err)
		return
	}
	patchReq.Header.Set("Authorization", "Bearer "+tokenStr)
	patchReq.Header.Set("Content-Type", "application/json")

	patchResp, err := client.Do(patchReq)
	if err != nil {
		fmt.Println("Error making PATCH request:", err)
		return
	}
	defer patchResp.Body.Close()

	patchBody, err := io.ReadAll(patchResp.Body)
	if err != nil {
		fmt.Println("Error reading PATCH response:", err)
		return
	}

	var patchResponse map[string]interface{}
	if err := json.Unmarshal(patchBody, &patchResponse); err != nil {
		fmt.Println("Error parsing PATCH response:", err)
		return
	}

	if patchResp.StatusCode != 200 {
		fmt.Println("PATCH API Error:", patchResponse["error"])
		return
	}

}
