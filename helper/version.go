package helper

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// UpdateCurrentVersion updates the current.json file with the new version
func UpdateCurrentVersion(cloudtmDir, version string) error {
	currentFile := filepath.Join(cloudtmDir, "current.json")

	currentData := map[string]string{
		"current": version,
	}

	currentJSON, err := json.MarshalIndent(currentData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(currentFile, currentJSON, 0644)
}

// GetCurrentVersion reads the current version from current.json
func GetCurrentVersion(cloudtmDir string) (string, error) {
	currentFile := filepath.Join(cloudtmDir, "current.json")

	data, err := os.ReadFile(currentFile)
	if err != nil {
		return "", err
	}

	var currentData map[string]string
	if err := json.Unmarshal(data, &currentData); err != nil {
		return "", err
	}

	return currentData["current"], nil
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
