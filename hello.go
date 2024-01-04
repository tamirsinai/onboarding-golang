package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type File struct {
	name string
	size int
}

type Scan struct {
	total int
	files []File
}

func main() {
	fmt.Print("Enter repo url: ")
	var repoUrl string
	fmt.Scanln(&repoUrl)

	fmt.Print("Enter min file size: ")
	var fileSize int
	fmt.Scanln(&fileSize)

	if gitClone(repoUrl) != nil {
		return
	}

	scan := scanRepoFiles("onboarding-golang", fileSize)
	fmt.Printf("%+v\n", scan)

	if deleteRepoDir() != nil {
		return
	}
}

func gitClone(repoUrl string) error {
	cmd := exec.Command("git", "clone", repoUrl)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
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
