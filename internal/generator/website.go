package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/open2b/scriggo"
	"github.com/open2b/scriggo/native"

	"github.com/odinnordico/odinnordico.github.io/internal/logger"
	"github.com/odinnordico/odinnordico.github.io/internal/models"
	"github.com/odinnordico/odinnordico.github.io/internal/utils"
)

// WebsiteGenerator handles static website generation
type WebsiteGenerator struct {
	templatesDir string
	theme        string
	assetsDir    string
}

// NewWebsiteGenerator creates a new website generator
func NewWebsiteGenerator(templatesDir, theme, assetsDir string) *WebsiteGenerator {
	return &WebsiteGenerator{
		templatesDir: templatesDir,
		theme:        theme,
		assetsDir:    assetsDir,
	}
}

// Generate generates the complete static website
func (wg *WebsiteGenerator) Generate(data *models.ResumeData, outputDir, lang string, copyAssets bool) error {
	logger.Logger().Info("Generating static website...")

	// Clear and create output directory
	if err := wg.clearOutputDir(outputDir); err != nil {
		return fmt.Errorf("failed to clear and create output directory: %w", err)
	}

	// Generate main pages
	if err := wg.generateIndexPage(data, outputDir, lang); err != nil {
		return fmt.Errorf("failed to generate index page: %w", err)
	}

	// Copy static assets
	if copyAssets {
		if err := wg.copyAssets(outputDir); err != nil {
			return fmt.Errorf("failed to copy assets: %w", err)
		}
	}

	logger.Logger().Info("Website generation completed")
	return nil
}

// generateIndexPage generates the main index.html page
func (wg *WebsiteGenerator) generateIndexPage(data *models.ResumeData, outputDir, lang string) error {
	templatePath := filepath.Join(wg.templatesDir, wg.theme, "index.html.tmpl")

	// Read template file
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Define custom template functions and data
	globals := native.Declarations{
		"Data":        data,
		"Lang":        lang,
		"DefaultLang": utils.DefaultLang,
		"seq": func(n int) []int {
			seq := make([]int, n)
			for i := 0; i < n; i++ {
				seq[i] = i + 1
			}
			return seq
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}

	opts := &scriggo.BuildOptions{
		Globals: globals,
	}

	fs := scriggo.Files{
		"index.html": templateContent,
	}

	tmpl, err := scriggo.BuildTemplate(fs, "index.html", opts)
	if err != nil {
		return fmt.Errorf("failed to build template: %w", err)
	}

	outputPath := filepath.Join(outputDir, "index.html")
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Run(file, nil, nil); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	logger.Logger().Info("Generated index", "outputPath", outputPath)
	return nil
}

// copyAssets copies static assets to the output directory
func (wg *WebsiteGenerator) copyAssets(outputDir string) error {
	assetsOutputDir := filepath.Join(outputDir, "assets")
	if err := os.MkdirAll(assetsOutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create assets directory: %w", err)
	}

	// Copy CSS, JS, images, etc.
	excludedFiles := map[string]bool{
		"files/.gitkeep": true,
		"files/cv.pdf":   true,
	}

	return filepath.Walk(wg.assetsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(wg.assetsDir, path)
		if err != nil {
			return err
		}

		if excludedFiles[relPath] {
			return nil
		}

		outputPath := filepath.Join(assetsOutputDir, relPath)
		outputDir := filepath.Dir(outputPath)

		logger.Logger().Debug("Copying asset", "path", path, "outputPath", outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}

		return copyFile(path, outputPath)
	})
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

// clearOutputDir removes the output directory and recreates it
func (wg *WebsiteGenerator) clearOutputDir(outputDir string) error {
	logger.Logger().Info("Clearing output directory", "dir", outputDir)
	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("failed to remove output directory: %w", err)
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	return nil
}
