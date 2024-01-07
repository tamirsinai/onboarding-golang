package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Input struct {
	Clone_url string
	Size      int
}

type File struct {
	name string
	size int
}

type Scan struct {
	total int
	files []File
}

func main() {
	input, err := readInputFile()
	if err != nil {
		fmt.Println("Error with read input file:", err)
		return
	}

	if err := gitClone(input.Clone_url); err != nil {
		fmt.Println("Error with clone repo:", err)
		return
	}

	scan := scanRepoFiles("/tmp/cloned-projects", input.Size)
	fmt.Printf("%+v\n", scan)

	if err := deleteRepoDir(); err != nil {
		fmt.Println("Error with delete repo:", err)
		return
	}
}

func readInputFile() (Input, error) {
	jsonData, err := os.ReadFile("input.json")
	if err != nil {
		return Input{}, err
	}

	var input Input
	err = json.Unmarshal(jsonData, &input)
	if err != nil {
		return input, err
	}

	return input, err
}

func gitClone(repoUrl string) error {
	os.Mkdir("/tmp/cloned-projects", 0755)
	os.Chdir("/tmp/cloned-projects")
	cmd := exec.Command("git", "clone", repoUrl)
	err := cmd.Run()
	return err
}

func scanRepoFiles(root string, fileSize int) Scan {
	scan := Scan{
		total: 0,
		files: []File{},
	}

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if !info.IsDir() && !strings.Contains(path, "/.git/") && int(info.Size()) > fileSize {
			scan.total = scan.total + 1
			scan.files = append(scan.files, File{path, int(info.Size())})
		}

		return nil
	})

	return scan
}

func deleteRepoDir() error {
	err := os.RemoveAll("/tmp/cloned-projects")

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Directory deleted successfully.")

	return err
}
