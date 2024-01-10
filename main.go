package main

import (
	"github.com/tamirsinai/onboarding-golang/handlers/files"
	"github.com/tamirsinai/onboarding-golang/handlers/repos"
	"path/filepath"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	input, err := files.ReadInputFile()
	if err != nil {
		logger.Error("Error read input file:", zap.Error(err))
		return
	}

	if err := repos.CloneRepositoryToScan(input.CloneUrl); err != nil {
		logger.Error("Error clone repo:", zap.Error(err))
		return
	}

	scan, err := repos.ScanRepoFiles(repos.ClonedProjectsDir, input.Size)
	if err != nil {
		logger.Error("Error scanning repo files:", zap.Error(err))
		return
	}

	if err := files.WriteOutputFile(scan); err != nil {
		logger.Error("Error write output file:", zap.Error(err))
		return
	}
	
	path, err := filepath.Abs(files.OutputFileName)
	if err != nil {
		logger.Error("Error getting absolute path:", zap.Error(err))
		return
	}
	logger.Info(path)

	if err := repos.DeleteClonedProjectsDir(); err != nil {
		logger.Error("Error with delete repo:", zap.Error(err))
		return
	}
}
