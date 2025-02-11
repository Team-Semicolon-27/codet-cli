package funcs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Clone(codatLink string) {
	prefix := "http://localhost:3000/codat/"

	if !strings.HasPrefix(codatLink, prefix) {
		fmt.Println("Invalid codat link")
		return
	}

	codatId := strings.TrimPrefix(codatLink, prefix)
	url := fmt.Sprintf("http://localhost:3000/api/codat/clone/%s", codatId)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching file:", err)
		return
	}
	defer resp.Body.Close()

	contentDisposition := resp.Header.Get("Content-Disposition")
	filename := "downloaded_codat.txt"

	if contentDisposition != "" {
		parts := strings.Split(contentDisposition, "filename=")
		if len(parts) > 1 {
			filename = strings.Trim(parts[1], `"`)
		}
	}

	dirName := strings.TrimSuffix(filename, filepath.Ext(filename))

	err = os.Mkdir(dirName, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating directory:", err)
		return
	}

	filePath := filepath.Join(dirName, filename)

	outFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Printf("File downloaded successfully as %s inside directory %s!\n", filename, dirName)
}
