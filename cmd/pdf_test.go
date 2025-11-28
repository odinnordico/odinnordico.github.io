package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/odinnordico/odinnordico.github.io/internal/models"
)

func TestDetectLanguages(t *testing.T) {
	// Create temp directory structure
	tempDir, err := os.MkdirTemp("", "test_detect_lang")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create lang directory with subdirectories
	langDir := filepath.Join(tempDir, "lang")
	if err := os.MkdirAll(filepath.Join(langDir, "es"), 0755); err != nil {
		t.Fatalf("Failed to create es dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(langDir, "fr"), 0755); err != nil {
		t.Fatalf("Failed to create fr dir: %v", err)
	}

	t.Run("Auto-detect all languages", func(t *testing.T) {
		langs := detectLanguages(tempDir, "")
		assert.Contains(t, langs, "en")
		assert.Contains(t, langs, "es")
		assert.Contains(t, langs, "fr")
	})

	t.Run("Specific language only", func(t *testing.T) {
		langs := detectLanguages(tempDir, "es")
		assert.Equal(t, []string{"es"}, langs)
	})

	t.Run("Default language", func(t *testing.T) {
		langs := detectLanguages(tempDir, "en")
		assert.Contains(t, langs, "en")
		// Should auto-detect others when default is specified
		assert.Contains(t, langs, "es")
		assert.Contains(t, langs, "fr")
	})
}

func TestGeneratePDF(t *testing.T) {
	// Create temp dirs
	tempDir, err := os.MkdirTemp("", "test_generate_pdf")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	outputDir := filepath.Join(tempDir, "output")
	templateDir := filepath.Join(tempDir, "templates", "default")

	// Create template structure
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("Failed to create template dir: %v", err)
	}

	// Create minimal template
	tmplContent := `
rows:
  - height: 10
    cols:
      - width: 12
        text:
          content: "Test Resume"
          size: 12
`
	if err := os.WriteFile(filepath.Join(templateDir, "resume.yaml.tmpl"), []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Change to temp dir so template path resolution works
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	data := &models.ResumeData{
		Basic: models.BasicData{
			Name: "Test User",
		},
	}

	t.Run("Generate PDF successfully", func(t *testing.T) {
		err := GeneratePDF(data, outputDir, "en", "default")
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(outputDir, "assets", "files", "resume.pdf"))
	})

	t.Run("Generate PDF with different language", func(t *testing.T) {
		err := GeneratePDF(data, outputDir, "es", "default")
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(outputDir, "assets", "files", "resume-es.pdf"))
	})
}
