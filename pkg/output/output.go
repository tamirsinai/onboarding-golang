package output

import (
	"github.com/tamirsinai/onboarding-golang/models"
	"os"
	"encoding/json"
)

const OutputFileName string = "output.json"

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