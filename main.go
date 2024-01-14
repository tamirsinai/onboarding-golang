package main

import (
	"github.com/tamirsinai/onboarding-golang/pkg/output"
	"github.com/tamirsinai/onboarding-golang/pkg/repo"
	"github.com/tamirsinai/onboarding-golang/pkg/input"
	"github.com/tamirsinai/onboarding-golang/pkg/scan"
	"path/filepath"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	input, err := input.ReadInputFile()
	if err != nil {
		logger.Error("Error read input file:", zap.Error(err))
		return
	}

	if err := repo.CloneRepositoryToScan(input.CloneUrl); err != nil {
		logger.Error("Error clone repo:", zap.Error(err))
		return
	}

	scan, err := scan.ScanRepoFiles(repo.ClonedProjectsDir, input.Size)
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

	if err := repo.DeleteClonedProjectsDir(); err != nil {
		logger.Error("Error with delete repo:", zap.Error(err))
		return
	}
}
