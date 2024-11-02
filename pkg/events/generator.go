package events

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// GenerateEvent creates a file with a UUID in the specified path
func GenerateEvent(path string) (string, error) {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create a file named with the current timestamp and a .json extension
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s.json", timestamp)
	filePath := filepath.Join(path, fileName)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Generate UUID
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to create random UUID: %w", err)
	}

	uuidData := map[string]string{"uuid": randomUUID.String()}
	jsonData, err := json.Marshal(uuidData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write data to file
	if _, err := file.Write(jsonData); err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}

	return filePath, nil
}
