package main

import (
	"github.com/tamirsinai/onboarding-golang/handlers/files"
	"github.com/tamirsinai/onboarding-golang/handlers/repos"
	"fmt"
	"path/filepath"
)

func main() {
	input, err := files.ReadInputFile()
	if err != nil {
		fmt.Println("Error with read input file:", err)
		return
	}

	if err := repos.CloneRepositoryToScan(input.CloneUrl); err != nil {
		fmt.Println("Error with clone repo:", err)
		return
	}

	scan := repos.ScanRepoFiles(repos.ClonedProjectsDir, input.Size)

	if err := files.WriteOutputFile(scan); err != nil {
		fmt.Println("Error with write output file:", err)
	}

	path, err := filepath.Abs("output.json")
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
	}
	fmt.Println(path)

	if err := repos.DeleteRepoDir(); err != nil {
		fmt.Println("Error with delete repo:", err)
		return
	}
}
