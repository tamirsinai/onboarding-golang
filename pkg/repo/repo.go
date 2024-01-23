package repo

import (
	"os"
	"gopkg.in/src-d/go-git.v4"
)

const ClonedProjectsDir string = "/tmp/cloned-projects"

func CloneRepositoryToScan(repoUrl string) error {
	if err := os.Mkdir(ClonedProjectsDir, 0755); err != nil {
		return err
	}

	_, err := git.PlainClone(ClonedProjectsDir, false, &git.CloneOptions{
        URL:      repoUrl,
        Progress: os.Stdout,
    })
	if err != nil {
        return err
    }
	return nil
	// cmd := exec.Command("git", "clone", repoUrl, ClonedProjectsDir)
	// output, err := cmd.CombinedOutput()
	// return errors.Wrap(err, string(output))
}

func DeleteClonedProjectsDir() error {
	return os.RemoveAll(ClonedProjectsDir)
}