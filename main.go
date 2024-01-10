package main

import (
	"github.com/tamirsinai/onboarding-golang/pkg/output"
	"github.com/tamirsinai/onboarding-golang/pkg/clone_repo"
	"github.com/tamirsinai/onboarding-golang/pkg/read_input"
	"github.com/tamirsinai/onboarding-golang/pkg/scan_repo_files"
	"github.com/tamirsinai/onboarding-golang/pkg/delete_cloned_projects_dir"
	"path/filepath"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	input, err := read_input.ReadInputFile()
	if err != nil {
		logger.Error("Error read input file:", zap.Error(err))
		return
	}

	if err := clone_repo.CloneRepositoryToScan(input.CloneUrl); err != nil {
		logger.Error("Error clone repo:", zap.Error(err))
		return
	}

	scan, err := scan_repo_files.ScanRepoFiles(clone_repo.ClonedProjectsDir, input.Size)
	if err != nil {
		logger.Error("Error scanning repo files:", zap.Error(err))
		return
	}

	if err := output.WriteOutputFile(scan); err != nil {
		logger.Error("Error write output file:", zap.Error(err))
		return
	}
	
	path, err := filepath.Abs(output.OutputFileName)
	if err != nil {
		logger.Error("Error getting absolute path:", zap.Error(err))
		return
	}
	logger.Info(path)

	if err := delete_cloned_projects_dir.DeleteClonedProjectsDir(); err != nil {
		logger.Error("Error with delete repo:", zap.Error(err))
		return
	}
}
