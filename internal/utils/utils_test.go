package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateDirectories(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_data")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Valid directory",
			path:    tempDir,
			wantErr: false,
		},
		{
			name:    "Non-existent directory",
			path:    filepath.Join(tempDir, "non_existent"),
			wantErr: true,
		},
		{
			name:    "Path is a file",
			path:    filepath.Join(tempDir, "test_file"),
			wantErr: true,
		},
	}

	// Create a dummy file for the "Path is a file" test case
	file, err := os.Create(filepath.Join(tempDir, "test_file"))
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDirectories(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDirectories() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
