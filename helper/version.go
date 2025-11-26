package helper

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// UpdateCurrentVersion updates the current.json file with the new version and status
func UpdateCurrentVersion(cloudtmDir, version string, status bool) error {
	currentFile := filepath.Join(cloudtmDir, "current.json")

	currentData := map[string]interface{}{
		"current": version,
		"status":  status,
	}

	currentJSON, err := json.MarshalIndent(currentData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(currentFile, currentJSON, 0644)
}

// GetCurrentVersion reads the current version and status from current.json
func GetCurrentVersion(cloudtmDir string) (string, bool, error) {
	currentFile := filepath.Join(cloudtmDir, "current.json")

	data, err := os.ReadFile(currentFile)
	if err != nil {
		return "", false, err
	}

	var currentData map[string]interface{}
	if err := json.Unmarshal(data, &currentData); err != nil {
		return "", false, err
	}

	version, _ := currentData["current"].(string)
	status, _ := currentData["status"].(bool)

	return version, status, nil
}

// SetCurrentStatus updates only the status field in current.json
func SetCurrentStatus(cloudtmDir string, status bool) error {
	currentFile := filepath.Join(cloudtmDir, "current.json")

	// Read existing data
	data, err := os.ReadFile(currentFile)
	if err != nil {
		return err
	}

	var currentData map[string]interface{}
	if err := json.Unmarshal(data, &currentData); err != nil {
		return err
	}

	// Update only status
	currentData["status"] = status

	// Write back
	currentJSON, err := json.MarshalIndent(currentData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(currentFile, currentJSON, 0644)
}

// UpdateRollbackVersion updates the rollback.json file with the rollback version
func UpdateRollbackVersion(cloudtmDir, version string) error {
	rollbackFile := filepath.Join(cloudtmDir, "rollback.json")

	rollbackData := map[string]string{
		"rollback": version,
	}

	rollbackJSON, err := json.MarshalIndent(rollbackData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(rollbackFile, rollbackJSON, 0644)
}

// GetRollbackVersion reads the rollback version from rollback.json
func GetRollbackVersion(cloudtmDir string) (string, error) {
	rollbackFile := filepath.Join(cloudtmDir, "rollback.json")

	data, err := os.ReadFile(rollbackFile)
	if err != nil {
		return "", err
	}

	var rollbackData map[string]string
	if err := json.Unmarshal(data, &rollbackData); err != nil {
		return "", err
	}

	return rollbackData["rollback"], nil
}

// IsRollbackEmpty checks if rollback.json has empty rollback field
func IsRollbackEmpty(cloudtmDir string) (bool, error) {
	version, err := GetRollbackVersion(cloudtmDir)
	if err != nil {
		return false, err
	}
	return version == "", nil
}
