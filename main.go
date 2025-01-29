package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

const repoDir = ".codet"

type Commit struct {
	Hash       string
	Message    string
	Timestamp  time.Time
	Deltas     map[string]string
	ParentHash string
}

func hashContent(content []byte) string {
	h := sha1.New()
	h.Write(content)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func calculateDelta(oldContent, newContent string) string {
	if oldContent == newContent {
		return ""	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldContent, newContent, false)
	return dmp.DiffToDelta(diffs)
}


func applyDelta(baseContent, delta string) (string, error) {
	dmp := diffmatchpatch.New()
	diffs, err := dmp.DiffFromDelta(baseContent, delta)
	if err != nil {
		return "", fmt.Errorf("error parsing delta: %v", err)
	}

	patches := dmp.PatchMake(diffs)
	appliedContent, results := dmp.PatchApply(patches, baseContent)

	if !results[0] {
		return "", fmt.Errorf("error applying patch")
	}

	return appliedContent, nil
}

func loadIndex() map[string]string {
	data, err := os.ReadFile(filepath.Join(repoDir, "index"))
	if err != nil {
		fmt.Printf("Error reading index file: %v\n", err)
		return nil 
	}
	
	var index map[string]string
	if err := json.Unmarshal(data, &index); err != nil {
		fmt.Printf("Error unmarshalling index data: %v\n", err)
		return nil
	}
	return index
}

func saveIndex(index map[string]string) {
	data, _ := json.Marshal(index)
	ioutil.WriteFile(filepath.Join(repoDir, "index"), data, 0644)
}

func loadHead() (string, error) {
	data, err := os.ReadFile(filepath.Join(repoDir, "HEAD"))
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %w", err)
	}
	return string(data), nil
}

func saveHead(hash string) error {
	err := os.WriteFile(filepath.Join(repoDir, "HEAD"), []byte(hash), 0644)
	if err != nil {
		return fmt.Errorf("failed to write HEAD: %w", err)
	}
	return nil
}


func loadCommits() []Commit {
	data, err := ioutil.ReadFile(filepath.Join(repoDir, "commits"))
	if err != nil {
		panic(err)
	}
	var commits []Commit
	json.Unmarshal(data, &commits)
	return commits
}

func saveCommits(commits []Commit) {
	data, _ := json.Marshal(commits)
	ioutil.WriteFile(filepath.Join(repoDir, "commits"), data, 0644)
}

func initRepo(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		os.Mkdir(repoDir, 0755)
		ioutil.WriteFile(filepath.Join(repoDir, "HEAD"), []byte(""), 0644)
		ioutil.WriteFile(filepath.Join(repoDir, "index"), []byte("{}"), 0644)
		ioutil.WriteFile(filepath.Join(repoDir, "commits"), []byte("[]"), 0644)
		fmt.Println("Initialized empty repository in", repoDir)
	} else {
		fmt.Println("Repository already initialized.")
	}
}

func addFile(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: mygit add <filename>")
		return
	}
	filename := args[0]

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	fileHash := hashContent(content)

	index := loadIndex()
	index[filename] = fileHash

	saveIndex(index)
	fmt.Println("Added file:", filename)
}

func commit(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: mygit commit <message>")
		return
	}
	message := args[0]

	headHash, err := loadHead()
	if err != nil {
		fmt.Println("Error loading HEAD:", err)
		return 
	}
	index := loadIndex()
	commits := loadCommits()

	newDeltas := make(map[string]string)

	for filename, _ := range index {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filename, err)
			continue
		}

		var baseContent string
		if headHash != "" {
			baseContent = getFileContentFromCommit(headHash, filename, commits)
		}
		delta := calculateDelta(baseContent, string(content))
		newDeltas[filename] = delta
	}

	newCommit := Commit{
		Hash:       hashContent([]byte(message + time.Now().String())),
		Message:    message,
		Timestamp:  time.Now(),
		Deltas:     newDeltas,
		ParentHash: headHash,
	}
	commits = append(commits, newCommit)

	saveCommits(commits)
	saveHead(newCommit.Hash)
	fmt.Println("Committed with message:", message)
}

func status(cmd *cobra.Command, args []string) {
	index := loadIndex()

	modifiedFiles := []string{}
	missingFiles := []string{}
	unchangedFiles := []string{}

	for filename, hash := range index {
		content, err := os.ReadFile(filename)
		if err != nil {
			missingFiles = append(missingFiles, filename)
			continue
		}

		currentHash := hashContent(content)
		if currentHash != hash {
			modifiedFiles = append(modifiedFiles, filename)
		} else {
			unchangedFiles = append(unchangedFiles, filename)
		}
	}

	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	untrackedFiles := []string{}
	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			if _, tracked := index[filename]; !tracked {
				untrackedFiles = append(untrackedFiles, filename)
			}
		}
	}

	fmt.Println("Status:")
	if len(modifiedFiles) > 0 {
		fmt.Println("Modified files:")
		for _, file := range modifiedFiles {
			fmt.Println("  -", file)
		}
	}
	if len(missingFiles) > 0 {
		fmt.Println("Missing files:")
		for _, file := range missingFiles {
			fmt.Println("  -", file)
		}
	}
	if len(untrackedFiles) > 0 {
		fmt.Println("Untracked files:")
		for _, file := range untrackedFiles {
			fmt.Println("  -", file)
		}
	}
	if len(modifiedFiles) == 0 && len(missingFiles) == 0 && len(untrackedFiles) == 0 {
		fmt.Println("No changes.")
	}
}


func getFileContentFromCommit(commitHash, filename string, commits []Commit) string {
	for _, commit := range commits {
		if commit.Hash == commitHash {
			delta, ok := commit.Deltas[filename]
			if !ok {
				return ""
			}
			baseContent := getFileContentFromCommit(commit.ParentHash, filename, commits)

			appliedContent, err := applyDelta(baseContent, delta)
			if err != nil {
				fmt.Println("Error applying delta:", err)
				return ""
			}

			return appliedContent
		}
	}
	return ""
}


func main() {
	var rootCmd = &cobra.Command{
		Use:   "codet",
		Short: "A cli tool for codet",
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new repository",
		Run:   initRepo,
	}

	var addCmd = &cobra.Command{
		Use:   "add <filename>",
		Short: "Add a file to the staging area",
		Run:   addFile,
	}

	var commitCmd = &cobra.Command{
		Use:   "commit <message>",
		Short: "Commit changes with a message",
		Run:   commit,
	}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show the status of tracked files",
		Run:   status,
	}

	rootCmd.AddCommand(initCmd, addCmd, commitCmd, statusCmd)
	rootCmd.Execute()
}
