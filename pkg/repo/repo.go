package repo

import (
	"os"
	"os/exec"
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

func DeleteClonedProjectsDir() error {
	return os.RemoveAll(ClonedProjectsDir)
}