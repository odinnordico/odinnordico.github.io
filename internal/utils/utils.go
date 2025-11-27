// Package utils provides utility functions for the odinnordico.github.io application.
package utils

import (
	"fmt"
	"os"
)

const DefaultLang = "en"

// ValidateDirectories checks if the specified data directory exists and is accessible.
// Returns an error if the directory doesn't exist or is not a directory.
func ValidateDirectories(dataDir string) error {
	dataInfo, err := os.Stat(dataDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("data directory not found: %s", dataDir)
	}
	if err != nil {
		return fmt.Errorf("cannot access data directory: %w", err)
	}
	if !dataInfo.IsDir() {
		return fmt.Errorf("data path is not a directory: %s", dataDir)
	}

	return nil
}
