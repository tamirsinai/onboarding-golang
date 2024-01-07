package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const clonedProjectsDir string = "/tmp/cloned-projects"
const inputFileName string = "input.json"

type Input struct {
	CloneUrl string `json:"clone_url"`
	Size     int
}

type File struct {
	Name string
	Size int
}

type Scan struct {
	Total int
	Files []File
}

func main() {
	input, err := readInputFile()
	if err != nil {
		fmt.Println("Error with read input file:", err)
		return
	}

	if err := cloneRepositoryToScan(input.CloneUrl); err != nil {
		fmt.Println("Error with clone repo:", err)
		return
	}

	scan := scanRepoFiles(clonedProjectsDir, input.Size)

	if err := writeOutputFile(scan); err != nil {
		fmt.Println("Error with write output file:", err)
	}

	path, err := filepath.Abs("output.json")
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
	}
	fmt.Println(path)

	if err := deleteRepoDir(); err != nil {
		fmt.Println("Error with delete repo:", err)
		return
	}
}

func readInputFile() (Input, error) {
	jsonData, err := os.ReadFile(inputFileName)
	if err != nil {
		return Input{}, err
	}

	var input Input
	if err := json.Unmarshal(jsonData, &input); err != nil {
		return input, err
	}

	return input, err
}

func cloneRepositoryToScan(repoUrl string) error {
	if err := os.Mkdir(clonedProjectsDir, 0755); err != nil {
		fmt.Println("Error make dir:", err)
		return err
	}

	cmd := exec.Command("git", "clone", repoUrl, clonedProjectsDir)
	err := cmd.Run()
	return err
}

func scanRepoFiles(repoPath string, fileSizeLimit int) Scan {
	scan := Scan{}

	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if !info.IsDir() && !strings.Contains(path, "/.git/") && int(info.Size()) > fileSizeLimit {
			scan.Total = scan.Total + 1
			scan.Files = append(scan.Files, File{path, int(info.Size())})
		}

		return nil
	})

	return scan
}

func writeOutputFile(scan Scan) error {
	jsonData, err := json.Marshal(scan)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	err = os.WriteFile("output.json", jsonData, 0644)
	return err
}

func deleteRepoDir() error {
	return os.RemoveAll(clonedProjectsDir)
}
