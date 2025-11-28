package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/odinnordico/odinnordico.github.io/internal/models"
)

func TestNewWebsiteGenerator(t *testing.T) {
	wg := NewWebsiteGenerator("templates", "default", "assets")
	assert.NotNil(t, wg)
}

func TestWebsiteGenerator_Generate(t *testing.T) {
	// Create temp dirs
	tempDir, err := os.MkdirTemp("", "test_website_gen")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	outputDir := filepath.Join(tempDir, "output")
	templatesDir := filepath.Join(tempDir, "templates")
	assetsDir := filepath.Join(tempDir, "assets")
	theme := "default"

	// Create template structure
	if err := os.MkdirAll(filepath.Join(templatesDir, theme), 0755); err != nil {
		t.Fatalf("Failed to create template dir structure: %v", err)
	}
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		t.Fatalf("Failed to create assets dir: %v", err)
	}

	// Create dummy index.html.tmpl
	tmplContent := `<html><body><h1>{{ Data.Basic.Name }}</h1></body></html>`
	if err := os.WriteFile(filepath.Join(templatesDir, theme, "index.html.tmpl"), []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Create dummy asset
	if err := os.WriteFile(filepath.Join(assetsDir, "style.css"), []byte("body {}"), 0644); err != nil {
		t.Fatalf("Failed to create asset file: %v", err)
	}

	wg := NewWebsiteGenerator(templatesDir, theme, assetsDir)
	data := &models.ResumeData{
		Basic: models.BasicData{
			Name: "John Doe",
		},
	}

	t.Run("Generate website successfully", func(t *testing.T) {
		err := wg.Generate(data, outputDir, "en", true)
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(outputDir, "index.html"))
		assert.FileExists(t, filepath.Join(outputDir, "assets", "style.css"))
	})
}
