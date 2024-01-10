package files

import (
	"github.com/tamirsinai/onboarding-golang/models"
	"os"
	"encoding/json"
)

const inputFileName string = "input.json"
const OutputFileName string = "output.json"

func ReadInputFile() (*models.Input, error) {
	jsonData, err := os.ReadFile(inputFileName)
	if err != nil {
		return nil, err
	}

	var input models.Input
	if err := json.Unmarshal(jsonData, &input); err != nil {
		return nil, err
	}

	return &input, nil
}

func WriteOutputFile(scan *models.Scan) error {
	jsonData, err := json.Marshal(&scan)
	if err != nil {
		return err
	}
	if err := os.WriteFile(OutputFileName, jsonData, 0644); err != nil {
		return err
	}
	return nil
}