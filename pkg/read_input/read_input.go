package read_input

import (
	"github.com/tamirsinai/onboarding-golang/models"
	"os"
	"encoding/json"
)

const inputFileName string = "input.json"

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