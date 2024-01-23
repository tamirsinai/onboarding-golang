package scan

import (
	"os"
	"path/filepath"
	"github.com/tamirsinai/onboarding-golang/models"
)

func ScanRepoFiles(repoPath string, fileSizeLimit int) (*models.Scan, error) {
	scan := models.Scan{}

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == ".git" {
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