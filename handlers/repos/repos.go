package repos

import (
	"os"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"github.com/tamirsinai/onboarding-golang/modules"
)

const ClonedProjectsDir string = "/tmp/cloned-projects"

func CloneRepositoryToScan(repoUrl string) error {
	if err := os.Mkdir(ClonedProjectsDir, 0755); err != nil {
		fmt.Println("Error make dir:", err)
		return err
	}

	cmd := exec.Command("git", "clone", repoUrl, ClonedProjectsDir)
	err := cmd.Run()
	return err
}

func ScanRepoFiles(repoPath string, fileSizeLimit int) modules.Scan {
	scan := modules.Scan{}

	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if !info.IsDir() && !strings.Contains(path, "/.git/") && int(info.Size()) > fileSizeLimit {
			scan.Total = scan.Total + 1
			scan.Files = append(scan.Files, modules.File{path, int(info.Size())})
		}

		return nil
	})

	return scan
}

func DeleteRepoDir() error {
	return os.RemoveAll(ClonedProjectsDir)
}