package delete_cloned_projects_dir

import (
	"os"
	"github.com/tamirsinai/onboarding-golang/pkg/clone_repo"
)

func DeleteClonedProjectsDir() error {
	return os.RemoveAll(clone_repo.ClonedProjectsDir)
}