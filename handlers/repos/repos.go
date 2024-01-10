package repos

import (
	"os"
	"os/exec"
	"path/filepath"
	"github.com/tamirsinai/onboarding-golang/models"
	"github.com/pkg/errors"
)

const ClonedProjectsDir string = "/tmp/cloned-projects"

func CloneRepositoryToScan(repoUrl string) error {
	if err := os.Mkdir(ClonedProjectsDir, 0755); err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", repoUrl, ClonedProjectsDir)
	output, err := cmd.CombinedOutput()
	return errors.Wrap(err, string(output))
}

func ScanRepoFiles(repoPath string, fileSizeLimit int) (*models.Scan, error) {
	subDirToSkip := ".git"
	scan := models.Scan{}

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == subDirToSkip {
			return filepath.SkipDir
		}

		if !info.IsDir() && int(info.Size()) > fileSizeLimit {
			scan.Total = scan.Total + 1
			scan.Files = append(scan.Files, models.File{Name: path, Size: int(info.Size())})

		}
		return nil
	})

	return &scan, err
}

func DeleteClonedProjectsDir() error {
	return os.RemoveAll(ClonedProjectsDir)
}