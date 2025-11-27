// Package generator provides functionality for generating PDF and website outputs from resume data.
// It includes a template engine for parsing YAML-based templates and rendering them into PDFs.
package generator

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
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

// RenderTemplate renders the parsed template into the Maroto PDF instance.
// It registers the footer (if present) and adds all rows to the document.
func RenderTemplate(m core.Maroto, t *Template) {
	if t.Footer != nil {
		cols := buildColumns(t.Footer.Cols)
		m.RegisterFooter(row.New(t.Footer.Height).Add(cols...))
	}

	for _, r := range t.Rows {
		cols := buildColumns(r.Cols)
		m.AddRow(r.Height, cols...)
	}
}

// buildColumns converts template column definitions into Maroto columns.
func buildColumns(cols []Col) []core.Col {
	result := make([]core.Col, 0, len(cols))
	for _, col := range cols {
		switch {
		case col.Text != nil:
			result = append(result, createTextCol(col.Width, col.Text))
		case col.Line != nil:
			result = append(result, createLineCol(col.Width, col.Line))
		case col.Image != nil:
			result = append(result, createImageCol(col.Width, col.Image))
		default:
			// Empty column
			result = append(result, text.NewCol(col.Width, ""))
		}
	}
	return result
}

// createTextCol creates a text column with the specified properties.
func createTextCol(width int, p *TextProp) core.Col {
	prop := props.Text{
		Size:  float64(p.Size),
		Style: parseStyle(p.Style),
		Align: parseAlign(p.Align),
	}

	if p.Color != nil {
		prop.Color = &props.Color{
			Red:   p.Color.Red,
			Green: p.Color.Green,
			Blue:  p.Color.Blue,
		}
	}

	if p.Hyperlink != "" {
		prop.Hyperlink = &p.Hyperlink
	}

	return text.NewCol(width, p.Content, prop)
}

// createLineCol creates a line column with the specified properties.
func createLineCol(width int, p *LineProp) core.Col {
	prop := props.Line{
		Thickness: p.Thickness,
	}

	if p.Color != nil {
		prop.Color = &props.Color{
			Red:   p.Color.Red,
			Green: p.Color.Green,
			Blue:  p.Color.Blue,
		}
	}

	return line.NewCol(width, prop)
}

// createImageCol creates an image column with the specified properties.
func createImageCol(width int, p *ImageProp) core.Col {
	rectProps := props.Rect{}
	if p.Percent > 0 {
		rectProps.Percent = p.Percent
	}
	if p.Center {
		rectProps.Center = p.Center
	}
	return image.NewFromFileCol(width, p.Path, rectProps)
}

// parseStyle converts a string style name to a fontstyle constant.
func parseStyle(style string) fontstyle.Type {
	switch strings.ToLower(style) {
	case "bold":
		return fontstyle.Bold
	case "italic":
		return fontstyle.Italic
	case "bolditalic":
		return fontstyle.BoldItalic
	default:
		return fontstyle.Normal
	}
}

// parseAlign converts a string alignment name to an align constant.
func parseAlign(alignment string) align.Type {
	switch strings.ToLower(alignment) {
	case "center":
		return align.Center
	case "right":
		return align.Right
	case "justify":
		return align.Justify
	default:
		return align.Left
	}
}
