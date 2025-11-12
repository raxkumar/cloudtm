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
