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
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldContent, newContent, false)
	return dmp.DiffToDelta(diffs)
}

func applyDelta(baseContent, delta string) string {
	dmp := diffmatchpatch.New()
	diffs, _ := dmp.DiffFromDelta(baseContent, delta)
	aplliedDelta, _ := dmp.PatchApply(dmp.PatchMake(diffs), baseContent)
	return aplliedDelta
}

func loadIndex() map[string]string {
	data, err := ioutil.ReadFile(filepath.Join(repoDir, "index"))
	if err != nil {
		panic(err)
	}
	var index map[string]string
	json.Unmarshal(data, &index)
	return index
}

func saveIndex(index map[string]string) {
	data, _ := json.Marshal(index)
	ioutil.WriteFile(filepath.Join(repoDir, "index"), data, 0644)
}

func loadHead() string {
	data, _ := ioutil.ReadFile(filepath.Join(repoDir, "HEAD"))
	return string(data)
}

func saveHead(hash string) {
	ioutil.WriteFile(filepath.Join(repoDir, "HEAD"), []byte(hash), 0644)
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

	headHash := loadHead()
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
	for filename, hash := range index {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("File missing: %s\n", filename)
			continue
		}
		currentHash := hashContent(content)
		if currentHash != hash {
			fmt.Printf("Modified: %s\n", filename)
		} else {
			fmt.Printf("Unchanged: %s\n", filename)
		}
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
			return applyDelta(baseContent, delta)
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
