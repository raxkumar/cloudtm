package helper

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// TerraformState represents the structure of terraform.tfstate
type TerraformState struct {
	Resources []interface{} `json:"resources"`
}

// IsStateEmpty checks if terraform.tfstate has empty resources array
func IsStateEmpty(workingDir string) (bool, error) {
	stateFile := filepath.Join(workingDir, "terraform.tfstate")

	// Check if state file exists
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		return false, err
	}

	// Read state file
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return false, err
	}

	// Parse JSON
	var state TerraformState
	if err := json.Unmarshal(data, &state); err != nil {
		return false, err
	}

	// Check if resources array is empty
	return len(state.Resources) == 0, nil
}
