// Package generator provides functionality for generating PDF and website outputs from resume data.
package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/grafana/gofpdf"

	"github.com/odinnordico/odinnordico.github.io/internal/logger"
	"github.com/odinnordico/odinnordico.github.io/internal/models"
	"github.com/odinnordico/odinnordico.github.io/internal/utils"
)

const (
	// PDF page margins in millimeters
	pdfMarginLeft   = 10.0
	pdfMarginTop    = 10.0
	pdfMarginRight  = 10.0
	pdfMarginBottom = 20.0 // Space for footer

	// A4 page dimensions and grid
	a4WidthMM       = 210.0
	a4UsableWidthMM = a4WidthMM - pdfMarginLeft - pdfMarginRight
	marotoCols      = 12.0 // Keep the 12-column grid concept
)

// PDFGenerator handles PDF resume generation using the gofpdf library.
type PDFGenerator struct {
	outputDir   string
	templateDir string
	theme       string
	tr          func(string) string
}

// NewPDFGenerator creates a new PDF generator with the specified configuration.
func NewPDFGenerator(outputDir, templateDir, theme string) (*PDFGenerator, error) {
	// Create a temporary PDF instance to get the translator
	// This is a bit of a workaround, but the translator is attached to the Fpdf struct in this library version
	tempPdf := gofpdf.New("P", "mm", "A4", "")
	tr := tempPdf.UnicodeTranslatorFromDescriptor("") // Default to cp1252

	return &PDFGenerator{
		outputDir:   outputDir,
		templateDir: templateDir,
		theme:       theme,
		tr:          tr,
	}, nil
}

// Generate creates a PDF resume from the provided resume data and language.
func (pg *PDFGenerator) Generate(data *models.ResumeData, lang string) error {
	if data == nil {
		return fmt.Errorf("resume data cannot be nil")
	}

	logger.Logger().Info("generating PDF resume with gofpdf")

	// Parse and execute template
	tmpl, err := pg.parseTemplate(data)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	// Configure PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(pdfMarginLeft, pdfMarginTop, pdfMarginRight)
	pdf.SetAutoPageBreak(true, pdfMarginBottom)
	pdf.AddPage()

	// Set default font
	pdf.SetFont("Arial", "", 11)

	// Render template into PDF
	pg.renderTemplate(pdf, tmpl)

	// Generate PDF document
	filename := "resume.pdf"
	if lang != utils.DefaultLang {
		filename = fmt.Sprintf("resume-%s.pdf", lang)
	}

	pdfPath := filepath.Join(pg.outputDir, "assets", "files", filename)
	if err := os.MkdirAll(filepath.Dir(pdfPath), 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	if err := pdf.OutputFileAndClose(pdfPath); err != nil {
		return fmt.Errorf("save PDF: %w", err)
	}

	logger.Logger().Info("PDF generated successfully", "file", pdfPath)
	return nil
}

// renderTemplate renders the parsed template into the gofpdf instance.
func (pg *PDFGenerator) renderTemplate(pdf *gofpdf.Fpdf, t *Template) {
	// Register footer if present
	if t.Footer != nil {
		pdf.SetFooterFunc(func() {
			// Save current position
			x, y := pdf.GetXY()

			// Position at 1.5 cm from bottom
			pdf.SetY(-15)

			// Render footer columns directly without checking for page breaks
			// We assume footer fits in the margin
			pg.renderCols(pdf, t.Footer.Height, t.Footer.Cols)

			// Restore position
			pdf.SetXY(x, y)
		})
	}

	for _, r := range t.Rows {
		pg.renderRow(pdf, &r)
	}
}

// renderRow renders a single row.
func (pg *PDFGenerator) renderRow(pdf *gofpdf.Fpdf, r *Row) {
	// Calculate max height for the row if not specified or if content exceeds it
	// For now, we respect the row height from YAML, but we might need to adjust based on content
	// gofpdf flows naturally, so we might just render columns.

	// However, to maintain the grid layout, we need to know the Y position.
	startY := pdf.GetY()

	// Check if we need a page break
	// Use pdfMarginBottom to ensure we don't overlap with footer
	if startY+r.Height > 297-pdfMarginBottom {
		pdf.AddPage()
		startY = pdf.GetY()
	}

	pg.renderCols(pdf, r.Height, r.Cols)

	// Move Y to the next row position
	// If the content was taller than r.Height, we should probably use that.
	// But for now, let's trust the template or the calculated max height.
	// We need to calculate the max height used by columns to be safe?
	// For now, use r.Height as before.
	pdf.SetY(startY + r.Height)
}

// renderCols renders the columns of a row at the current Y position.
func (pg *PDFGenerator) renderCols(pdf *gofpdf.Fpdf, rowHeight float64, cols []Col) {
	startY := pdf.GetY()
	currentX := pdfMarginLeft

	maxHeight := rowHeight

	for _, col := range cols {
		colWidth := (a4UsableWidthMM / marotoCols) * float64(col.Width)

		pdf.SetXY(currentX, startY)

		if col.Text != nil {
			h := pg.renderText(pdf, colWidth, col.Text)
			if h > maxHeight {
				maxHeight = h
			}
		} else if col.Line != nil {
			pg.renderLine(pdf, colWidth, rowHeight, col.Line)
		} else if col.Image != nil {
			pg.renderImage(pdf, colWidth, rowHeight, col.Image)
		}

		currentX += colWidth
	}
}

// renderText renders a text column.
func (pg *PDFGenerator) renderText(pdf *gofpdf.Fpdf, width float64, p *TextProp) float64 {
	// Set style
	style := ""
	if strings.Contains(strings.ToLower(p.Style), "bold") {
		style += "B"
	}
	if strings.Contains(strings.ToLower(p.Style), "italic") {
		style += "I"
	}
	pdf.SetFont("Arial", style, float64(p.Size))

	// Set color
	if p.Color != nil {
		pdf.SetTextColor(p.Color.Red, p.Color.Green, p.Color.Blue)
	} else {
		pdf.SetTextColor(0, 0, 0)
	}

	// Alignment
	align := "L"
	switch strings.ToLower(p.Align) {
	case "center":
		align = "C"
	case "right":
		align = "R"
	case "justify":
		align = "J"
	}

	// Hyperlink
	if p.Hyperlink != "" {
		// We can add a link annotation
		// gofpdf Link support is a bit manual.
		// For now, let's just write the text.
		// If we want clickable links, we need to use pdf.LinkString or similar.
		// But MultiCell doesn't easily support mixed links.
		// If the whole block is a link, we can do it.
		// Let's assume the whole block is a link if Hyperlink is set.
		linkID := pdf.AddLink()
		pdf.SetLink(linkID, 0, -1)
		pdf.LinkString(pdf.GetX(), pdf.GetY(), width, 5, p.Hyperlink) // Height is approximate
	}

	// Render text
	// MultiCell(w, h, txt, border, align, fill)
	// h is line height.
	// Maroto default line height is around 1.0-1.2 depending on font.
	// We use a slightly larger line height for better readability.
	lineHeight := float64(p.Size) * 0.3527 * 1.3

	// Translate text to handle special characters (e.g., tildes)
	content := pg.tr(p.Content)

	pdf.MultiCell(width, lineHeight, content, "", align, false)

	// Return the height used
	lines := pdf.SplitLines([]byte(content), width)
	return float64(len(lines)) * lineHeight
}

// renderLine renders a line column.
func (pg *PDFGenerator) renderLine(pdf *gofpdf.Fpdf, width, height float64, p *LineProp) {
	// Set color
	if p.Color != nil {
		pdf.SetDrawColor(p.Color.Red, p.Color.Green, p.Color.Blue)
	} else {
		pdf.SetDrawColor(200, 200, 200) // Default light gray
	}

	pdf.SetLineWidth(p.Thickness)

	// Draw line centered vertically in the row height?
	// Or at the top?
	// Maroto lines are usually centered or fill the width.
	// Let's draw a horizontal line in the middle of the height?
	// Or maybe it's a vertical separator?
	// Based on resume templates, it's usually a horizontal line.
	// Let's assume horizontal line at Y + height/2

	x, y := pdf.GetXY()
	lineY := y + (height / 2)

	pdf.Line(x, lineY, x+width, lineY)
}

// renderImage renders an image column.
func (pg *PDFGenerator) renderImage(pdf *gofpdf.Fpdf, width, height float64, p *ImageProp) {
	// Image(src, x, y, w, h, flow, tp, link, linkStr)
	// If h is 0, it auto-scales.

	x, y := pdf.GetXY()

	imgWidth := width
	if p.Percent > 0 {
		imgWidth = width * (p.Percent / 100.0)
	}

	// Center image if requested
	if p.Center {
		x += (width - imgWidth) / 2
	}

	// Check if file exists
	if _, err := os.Stat(p.Path); os.IsNotExist(err) {
		logger.Logger().Warn("Image not found", "path", p.Path)
		return
	}

	pdf.Image(p.Path, x, y, imgWidth, 0, false, "", 0, "")
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
		"calculateHeight":   calculateHeight, // We might not need this anymore but keep for template compatibility
		"getSocials":        getSocials,
		"chunkSocials":      chunkSocials,
		"assetPath":         assetPath,
	}
}

// Helper functions (kept from original)

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
// Kept for template compatibility, but logic might need adjustment for gofpdf if used in template.
func calculateHeight(text string, fontSize int, colWidth int) float64 {
	// This was tuned for Maroto. For gofpdf, it might be different.
	// But since it's used in the template to set row height, we should keep it or improve it.
	// Let's keep the logic for now.

	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}

	// Calculate column width in millimeters
	colWidthMM := (a4UsableWidthMM / marotoCols) * float64(colWidth)

	// Estimate character width
	// Maroto used 0.55 factor. gofpdf Arial might be similar.
	pointsToMM := 0.3527
	charWidthFactor := 0.55

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
	lineHeightFactor := 1.2 // gofpdf default is usually around 1.2
	lineHeightMM := float64(fontSize) * pointsToMM * lineHeightFactor
	totalHeight := lines * lineHeightMM

	return totalHeight + 1.0 // Padding
}
