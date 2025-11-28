package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/odinnordico/odinnordico.github.io/internal/models"
)

func TestNewPDFGenerator(t *testing.T) {
	pg, err := NewPDFGenerator("output", "templates", "default")
	assert.NoError(t, err)
	assert.NotNil(t, pg)
}

func TestPDFGenerator_Generate(t *testing.T) {
	// Create temp dirs
	tempDir, err := os.MkdirTemp("", "test_pdf_gen")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	outputDir := filepath.Join(tempDir, "output")
	templateDir := filepath.Join(tempDir, "templates")
	theme := "default"

	// Create template structure
	if err := os.MkdirAll(filepath.Join(templateDir, theme), 0755); err != nil {
		t.Fatalf("Failed to create template dir structure: %v", err)
	}

	// Create dummy template
	tmplContent := `
rows:
  - height: 10
    cols:
      - width: 12
        text:
          content: "Test Resume"
          size: 12
`
	if err := os.WriteFile(filepath.Join(templateDir, theme, "resume.yaml.tmpl"), []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	pg, err := NewPDFGenerator(outputDir, templateDir, theme)
	assert.NoError(t, err)

	data := &models.ResumeData{
		Basic: models.BasicData{
			Name: "John Doe",
		},
	}

	t.Run("Generate PDF successfully", func(t *testing.T) {
		err := pg.Generate(data, "en")
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(outputDir, "assets", "files", "resume.pdf"))
	})

	t.Run("Generate PDF with nil data", func(t *testing.T) {
		err := pg.Generate(nil, "en")
		assert.Error(t, err)
	})
}
