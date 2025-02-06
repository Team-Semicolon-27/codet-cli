package funcs

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

	outFile, err := os.Create(filename)
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

	fmt.Printf("File downloaded successfully as %s!\n", filename)
}
