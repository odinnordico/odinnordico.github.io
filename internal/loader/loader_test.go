package loader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadResumeData(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_data")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create dummy YAML files
	createYAMLFile(t, tempDir, "basic.yaml", "name: John Doe")
	createYAMLFile(t, tempDir, "professional.yaml", "title: Software Engineer")

	// Create language specific directory and file
	langDir := filepath.Join(tempDir, "lang", "es")
	if err := os.MkdirAll(langDir, 0755); err != nil {
		t.Fatalf("Failed to create lang dir: %v", err)
	}
	createYAMLFile(t, langDir, "basic.yaml", "name: Juan Perez")

	t.Run("Load default language", func(t *testing.T) {
		data, err := LoadResumeData(tempDir, "")
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", data.Basic.Name)
	})

	t.Run("Load specific language", func(t *testing.T) {
		data, err := LoadResumeData(tempDir, "es")
		assert.NoError(t, err)
		assert.Equal(t, "Juan Perez", data.Basic.Name)
	})

	t.Run("Load non-existent language", func(t *testing.T) {
		data, err := LoadResumeData(tempDir, "fr")
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", data.Basic.Name) // Should fallback to default or not merge anything
	})
}

func createYAMLFile(t *testing.T, dir, filename, content string) {
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create YAML file %s: %v", path, err)
	}
}

// Mocking models for test purpose if needed, but using actual models is better for integration-like unit test.
// Assuming models.ResumeData has Basic struct with Name field.
