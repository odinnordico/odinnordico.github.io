package generator

import (
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestParseTemplate(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_template")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a dummy template file
	tmplContent := `
rows:
  - height: 10
    cols:
      - width: 12
        text:
          content: "{{ .Name }}"
          size: 12
`
	tmplPath := filepath.Join(tempDir, "resume.yaml.tmpl")
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Test data
	data := struct {
		Name string
	}{
		Name: "John Doe",
	}

	// Test ParseTemplate
	t.Run("Valid template", func(t *testing.T) {
		tmpl, err := ParseTemplate(tmplPath, data, template.FuncMap{})
		assert.NoError(t, err)
		assert.NotNil(t, tmpl)
		assert.Len(t, tmpl.Rows, 1)
		assert.Equal(t, "John Doe", tmpl.Rows[0].Cols[0].Text.Content)
	})

	t.Run("Invalid template path", func(t *testing.T) {
		_, err := ParseTemplate("non_existent.tmpl", data, template.FuncMap{})
		assert.Error(t, err)
	})

	t.Run("Invalid template content", func(t *testing.T) {
		invalidTmplPath := filepath.Join(tempDir, "invalid.tmpl")
		if err := os.WriteFile(invalidTmplPath, []byte("{{ .Name }"), 0644); err != nil {
			t.Fatalf("Failed to create invalid template file: %v", err)
		}
		_, err := ParseTemplate(invalidTmplPath, data, template.FuncMap{})
		assert.Error(t, err)
	})
}
