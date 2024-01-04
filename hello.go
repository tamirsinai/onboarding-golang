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
	jsonData, err := os.ReadFile("input.json")
	if err != nil {
		fmt.Println("Error:", err)
	}
	var input Input
	err = json.Unmarshal(jsonData, &input)
	if err != nil {
		fmt.Println("Error:", err)
	}

	if gitClone(input.Clone_url) != nil {
		return
	}

	scan := scanRepoFiles("onboarding-golang", input.Size)
	fmt.Printf("%+v\n", scan)

	if deleteRepoDir() != nil {
		return
	}
}

func gitClone(repoUrl string) error {
	cmd := exec.Command("git", "clone", repoUrl)
	err := cmd.Run()

	if err != nil {
		fmt.Println("Error executing git clone:", err)
	} else {
		fmt.Println("Git clone successful.")
	}

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
	err := os.RemoveAll("onboarding-golang")

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Directory deleted successfully.")

	return err
}
