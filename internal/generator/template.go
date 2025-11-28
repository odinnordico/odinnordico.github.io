// Package generator provides functionality for generating PDF and website outputs from resume data.
// It includes a template engine for parsing YAML-based templates and rendering them into PDFs.
package generator

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Template represents the YAML structure of a resume template.
// It contains rows for the main content and an optional footer.
type Template struct {
	Rows   []Row `yaml:"rows"`
	Footer *Row  `yaml:"footer,omitempty"`
}

// Row represents a horizontal row in the PDF with a specified height and columns.
type Row struct {
	Height float64 `yaml:"height"`
	Cols   []Col   `yaml:"cols"`
}

// Col represents a column within a row, which can contain text, a line, or an image.
type Col struct {
	Width int        `yaml:"width"`
	Text  *TextProp  `yaml:"text,omitempty"`
	Line  *LineProp  `yaml:"line,omitempty"`
	Image *ImageProp `yaml:"image,omitempty"`
}

// TextProp defines properties for text content in a column.
type TextProp struct {
	Content   string `yaml:"content"`
	Size      int    `yaml:"size"`
	Style     string `yaml:"style"` // normal, bold, italic, bolditalic
	Align     string `yaml:"align"` // left, center, right, justify
	Color     *Color `yaml:"color,omitempty"`
	Hyperlink string `yaml:"hyperlink,omitempty"`
}

// LineProp defines properties for a horizontal line in a column.
type LineProp struct {
	Thickness float64 `yaml:"thickness"`
	Color     *Color  `yaml:"color,omitempty"`
}

// ImageProp defines properties for an image in a column.
type ImageProp struct {
	Path    string  `yaml:"path"`
	Percent float64 `yaml:"percent,omitempty"` // Size as percentage of column width
	Center  bool    `yaml:"center,omitempty"`  // Whether to center the image
}

// Color represents an RGB color value.
type Color struct {
	Red   int `yaml:"red"`
	Green int `yaml:"green"`
	Blue  int `yaml:"blue"`
}

// ParseTemplate reads a YAML template file, executes it as a Go template with the provided data,
// and parses the resulting YAML into a Template structure.
// Returns an error if the file cannot be read, the template cannot be parsed or executed,
// or the resulting YAML is invalid.
func ParseTemplate(path string, data interface{}, funcs template.FuncMap) (*Template, error) {
	tmplContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read template file: %w", err)
	}

	tmpl, err := template.New("resume").Funcs(funcs).Parse(string(tmplContent))
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	var t Template
	if err := yaml.Unmarshal(buf.Bytes(), &t); err != nil {
		// Log the generated YAML for debugging purposes if parsing fails
		fmt.Println("Generated YAML:\n", buf.String())
		return nil, fmt.Errorf("unmarshal YAML: %w", err)
	}

	return &t, nil
}
