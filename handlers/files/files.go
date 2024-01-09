package files

import (
	"example/modules"
	"os"
	"encoding/json"
	"fmt"
)

const inputFileName string = "input.json"

func ReadInputFile() (*modules.Input, error) {
	jsonData, err := os.ReadFile(inputFileName)
	if err != nil {
		return nil, err
	}

	var input modules.Input
	if err := json.Unmarshal(jsonData, &input); err != nil {
		return nil, err
	}

	return &input, err
}

func WriteOutputFile(scan modules.Scan) error {
	jsonData, err := json.Marshal(scan)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	err = os.WriteFile("output.json", jsonData, 0644)
	return err
}