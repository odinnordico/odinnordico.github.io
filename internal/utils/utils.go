// Package utils provides utility functions for the odinnordico.github.io application.
package utils

import (
	"fmt"
	"os"

	"github.com/odinnordico/odinnordico.github.io/internal/logger"
)

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

// GetLang extracts the two-letter language code from a potentially longer locale string.
// For example, "en_US.UTF-8" becomes "en".
// If the input is already 2 characters or less, it is returned unchanged.
func GetLang(lang string) string {
	logger.Logger().Debug("Extracting language code", "lang", lang)
	if len(lang) > 2 {
		return lang[:2]
	}
	logger.Logger().Debug("Language code extracted", "lang", lang)
	return lang
}
