package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/odinnordico/odinnordico.github.io/internal/models"
)

func TestGenerateWebsite(t *testing.T) {
	// Create temp dirs
	tempDir, err := os.MkdirTemp("", "test_generate_website")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	outputDir := filepath.Join(tempDir, "output")
	templateDir := filepath.Join(tempDir, "templates", "default")
	assetsDir := filepath.Join(tempDir, "assets")

	// Create directory structure
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		t.Fatalf("Failed to create assets dir: %v", err)
	}

	// Create minimal template
	tmplContent := `<html><body><h1>{{ Data.Basic.Name }}</h1></body></html>`
	if err := os.WriteFile(filepath.Join(templateDir, "index.html.tmpl"), []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Create dummy asset
	if err := os.WriteFile(filepath.Join(assetsDir, "style.css"), []byte("body {}"), 0644); err != nil {
		t.Fatalf("Failed to create asset: %v", err)
	}

	// Change to temp dir for path resolution
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	data := &models.ResumeData{
		Basic: models.BasicData{
			Name: "Test User",
		},
	}

	t.Run("Generate website with assets", func(t *testing.T) {
		err := GenerateWebsite(data, outputDir, "en", "default", true)
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(outputDir, "index.html"))
		assert.FileExists(t, filepath.Join(outputDir, "assets", "style.css"))
	})

	t.Run("Generate website without copying assets", func(t *testing.T) {
		outputDir2 := filepath.Join(tempDir, "output2")
		err := GenerateWebsite(data, outputDir2, "es", "default", false)
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(outputDir2, "index.html"))
		assert.NoDirExists(t, filepath.Join(outputDir2, "assets"))
	})
}
