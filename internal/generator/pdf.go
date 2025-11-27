// Package generator provides functionality for generating PDF and website outputs from resume data.
package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"github.com/odinnordico/odinnordico.github.io/internal/logger"
	"github.com/odinnordico/odinnordico.github.io/internal/models"
	"github.com/odinnordico/odinnordico.github.io/internal/utils"
)

const (
	// PDF page margins in millimeters
	pdfMarginLeft  = 10
	pdfMarginTop   = 10
	pdfMarginRight = 10

	// A4 page dimensions and grid
	a4WidthMM       = 210.0
	a4UsableWidthMM = a4WidthMM - pdfMarginLeft - pdfMarginRight
	marotoCols      = 12

	// Typography constants
	pointsToMM       = 0.3527 // 1 point = 0.3527 mm
	charWidthFactor  = 0.55   // Conservative estimate for variable width fonts
	lineHeightFactor = 1.1    // Line height multiplier for spacing
	heightPadding    = 0.5    // Additional padding in mm
)

// PDFGenerator handles PDF resume generation using the Maroto library.
type PDFGenerator struct {
	outputDir   string
	templateDir string
	theme       string
}

// NewPDFGenerator creates a new PDF generator with the specified configuration.
// It returns an error if the configuration is invalid (currently always returns nil).
func NewPDFGenerator(outputDir, templateDir, theme string) (*PDFGenerator, error) {
	return &PDFGenerator{
		outputDir:   outputDir,
		templateDir: templateDir,
		theme:       theme,
	}, nil
}

// Generate creates a PDF resume from the provided resume data and language.
// The PDF filename will include the language code if it's not English (e.g., "resume-es.pdf").
// Returns an error if template parsing, PDF generation, or file saving fails.
func (pg *PDFGenerator) Generate(data *models.ResumeData, lang string) error {
	if data == nil {
		return fmt.Errorf("resume data cannot be nil")
	}

	logger.Logger().Info("generating PDF resume with Maroto template")

	// Parse and execute template
	tmpl, err := pg.parseTemplate(data)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	// Configure PDF document
	cfg := config.NewBuilder().
		WithPageNumber(props.PageNumber{
			Pattern: "{current}",
			Place:   props.RightBottom,
		}).
		WithLeftMargin(pdfMarginLeft).
		WithTopMargin(pdfMarginTop).
		WithRightMargin(pdfMarginRight).
		Build()

	m := maroto.New(cfg)

	// Render template into PDF
	RenderTemplate(m, tmpl)

	// Generate PDF document
	document, err := m.Generate()
	if err != nil {
		return fmt.Errorf("generate PDF: %w", err)
	}

	// Save to file
	if err := pg.savePDF(document, lang); err != nil {
		return fmt.Errorf("save PDF: %w", err)
	}

	return nil
}

// parseTemplate loads and parses the YAML template with the resume data.
func (pg *PDFGenerator) parseTemplate(data *models.ResumeData) (*Template, error) {
	tmplPath := filepath.Join(pg.templateDir, pg.theme, "resume.yaml.tmpl")
	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("template file not found: %s", tmplPath)
	}

	funcs := pg.buildTemplateFuncs()
	tmpl, err := ParseTemplate(tmplPath, data, funcs)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// buildTemplateFuncs creates the template.FuncMap with all available template functions.
func (pg *PDFGenerator) buildTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// Data extraction
		"getEmail":   getEmail,
		"getPhone":   getPhone,
		"hasSocials": hasSocials,

		// Formatting
		"formatDate":        formatDate,
		"formatSkills":      formatSkills,
		"formatCurrentDate": formatCurrentDate,
		"escapeYAML":        escapeYAML,
		"splitLines":        splitLines,
		"lastURLPart":       lastURLPart,
		"calculateHeight":   calculateHeight,
		"getSocials":        getSocials,
		"chunkSocials":      chunkSocials,
		"assetPath":         assetPath,
	}
}

// savePDF saves the generated PDF document to the output directory.
func (pg *PDFGenerator) savePDF(document core.Document, lang string) error {
	filename := "resume.pdf"
	if lang != utils.DefaultLang {
		filename = fmt.Sprintf("resume-%s.pdf", lang)
	}

	pdfPath := filepath.Join(pg.outputDir, "assets", "files", filename)
	if err := os.MkdirAll(filepath.Dir(pdfPath), 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	if err := document.Save(pdfPath); err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	logger.Logger().Info("PDF generated successfully", "file", pdfPath)
	return nil
}

// Template helper functions

// getEmail extracts the email address from the resume data.
func getEmail(data *models.ResumeData) string {
	return getSocialValue(data, "Email")
}

// getPhone extracts the phone number from the resume data.
func getPhone(data *models.ResumeData) string {
	return getSocialValue(data, "Phone", "Mobile")
}

// getSocialValue extracts the URL of the first social entry matching one of the given names.
func getSocialValue(data *models.ResumeData, names ...string) string {
	for _, social := range data.Social {
		for _, name := range names {
			if strings.EqualFold(social.Name, name) {
				if strings.EqualFold(name, "Email") {
					return strings.TrimPrefix(social.URL, "mailto:")
				}
				return social.URL
			}
		}
	}
	return ""
}

// hasSocials checks if the resume has any social media links (excluding email and phone).
func hasSocials(data *models.ResumeData) bool {
	for _, s := range data.Social {
		if !isContactInfo(s.Name) {
			return true
		}
	}
	return false
}

// isContactInfo checks if a social entry is contact information (phone/mobile).
func isContactInfo(name string) bool {
	return strings.EqualFold(name, "Phone") ||
		strings.EqualFold(name, "Mobile")
}

// getSocials returns a list of social media links excluding contact info.
func getSocials(data *models.ResumeData) []models.Entity {
	var socials []models.Entity
	for _, s := range data.Social {
		if !isContactInfo(s.Name) {
			socials = append(socials, s)
		}
	}
	return socials
}

// chunkSocials splits a slice of entities into chunks of the specified size.
func chunkSocials(values []models.Entity, chunkSize int) [][]models.Entity {
	if chunkSize <= 0 {
		return nil
	}
	var chunks [][]models.Entity
	for i := 0; i < len(values); i += chunkSize {
		end := i + chunkSize
		if end > len(values) {
			end = len(values)
		}
		chunks = append(chunks, values[i:end])
	}
	return chunks
}

// formatDate formats a time
func formatDate(t *time.Time, format string) string {
	if t == nil {
		return "Present"
	}
	if format == "" {
		format = "Jan 2006"
	}
	if t.IsZero() {
		return ""
	}
	return t.Format(format)
}

// formatCurrentDate returns the current date formatted as specified.
func formatCurrentDate(format string) string {
	t := time.Now()
	return formatDate(&t, format)
}

// formatSkills formats skills as a comma-separated string.
func formatSkills(data *models.ResumeData) string {
	var skillStrings []string
	for _, s := range data.Skills {
		skillStrings = append(skillStrings, s.Name)
	}
	return strings.Join(skillStrings, ", ")
}

// escapeYAML escapes special characters for YAML content.
func escapeYAML(s string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"\"", "\\\"",
		"\n", "\\n",
	)
	return replacer.Replace(s)
}

// splitLines splits a string by newlines.
func splitLines(s string) []string {
	return strings.Split(s, "\n")
}

// lastURLPart extracts the last segment of a URL path.
func lastURLPart(url string) string {
	return lastPart(lastPart(url, "/"), ":")
}

// lastPart extracts the last segment of a string based on a separator.
func lastPart(url, separator string) string {
	parts := strings.Split(url, separator)
	return parts[len(parts)-1]
}

// assetPath converts a relative path to an absolute path based on the current working directory.
func assetPath(path string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, path)
}

// calculateHeight estimates the required height in millimeters for rendering text
// based on font size, column width, and text content.
func calculateHeight(text string, fontSize int, colWidth int) float64 {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}

	// Calculate column width in millimeters
	colWidthMM := (a4UsableWidthMM / marotoCols) * float64(colWidth)

	// Estimate character width
	charWidthMM := float64(fontSize) * pointsToMM * charWidthFactor

	// Calculate characters per line
	charsPerLine := int(colWidthMM / charWidthMM)
	if charsPerLine < 1 {
		charsPerLine = 1
	}

	// Count wrapped lines
	lines := 0.0
	textLines := strings.Split(text, "\n")

	for _, line := range textLines {
		lineLength := float64(len(line))
		if lineLength == 0 {
			lines++ // Empty line still takes space
			continue
		}

		wrappedLines := (lineLength + float64(charsPerLine) - 1) / float64(charsPerLine)
		lines += float64(int(wrappedLines))
	}

	// Calculate total height
	lineHeightMM := float64(fontSize) * pointsToMM * lineHeightFactor
	totalHeight := lines * lineHeightMM

	return totalHeight + heightPadding
}
