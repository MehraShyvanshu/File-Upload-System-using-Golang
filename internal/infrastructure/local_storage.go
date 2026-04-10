package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
)

// saves file locally
func SaveFileLocally(filename string, data []byte) (string, error) {

	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	// create uploads folder path
	uploadsDir := filepath.Join(cwd, "uploads")

	// create folder with proper permissions
	err = os.MkdirAll(uploadsDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create uploads directory: %w", err)
	}

	// create full file path
	filePath := filepath.Join(uploadsDir, filename)

	// write file
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// return relative path for storage in database
	relPath := filepath.Join("uploads", filename)
	return relPath, nil
}
