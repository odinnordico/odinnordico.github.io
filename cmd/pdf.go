// Package cmd provides command-line interface commands for the odinnordico.github.io application.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/odinnordico/odinnordico.github.io/internal/generator"
	"github.com/odinnordico/odinnordico.github.io/internal/loader"
	"github.com/odinnordico/odinnordico.github.io/internal/logger"
	"github.com/odinnordico/odinnordico.github.io/internal/models"
	"github.com/odinnordico/odinnordico.github.io/internal/utils"
)

// PdfCmd represents the pdf command for generating PDF resumes.
var PdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "Generate PDF resume from YAML resume data",
	Long:  `Generate creates a PDF resume from YAML data files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dataDir := viper.GetString("data-dir")
		theme := viper.GetString("theme")
		outputDir := viper.GetString("output-dir")

		logger.Logger().Info("Starting PDF generation...")
		logger.Logger().Info("Data directory", "dir", dataDir)
		logger.Logger().Info("Theme", "name", theme)

		// Validate input directory
		if err := utils.ValidateDirectories(dataDir); err != nil {
			return fmt.Errorf("directory validation failed: %w", err)
		}

		return GenerateMultiLanguagePdf(dataDir, outputDir, utils.DefaultLang, theme)
	},
}

func init() {
	PdfCmd.Flags().String("theme", "default", "Theme name to use")
	viper.BindPFlag("theme", PdfCmd.Flags().Lookup("theme"))
}

// detectLanguages determines which languages to generate PDFs for based on the target language
// and available language directories. If targetLang is empty or utils.DefaultLang, it auto-detects all available languages.
func detectLanguages(dataDir, targetLang string) []string {
	langDirName := "lang"
	languages := []string{utils.DefaultLang}

	if targetLang != "" && targetLang != utils.DefaultLang {
		return []string{targetLang}
	}

	// Auto-detect other languages from data/lang directory
	langDir := filepath.Join(dataDir, langDirName)
	filepath.WalkDir(langDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		logger.Logger().Debug("Language directory", "entry", d.Name())
		if d.IsDir() && d.Name() != utils.DefaultLang && d.Name() != langDirName {
			languages = append(languages, d.Name())
		}
		return nil
	})

	return languages
}

// GenerateMultiLanguagePdf generates a PDF resume for the specified language using the given data and theme.
// If targetLang is empty or utils.DefaultLang, it auto-detects all available languages.
func GenerateMultiLanguagePdf(dataDir, outputDir, targetLang, theme string) error {
	languages := detectLanguages(dataDir, targetLang)
	logger.Logger().Info("Target Lang", "lang", targetLang)
	logger.Logger().Info("Languages to generate", "langs", languages)

	// Generate PDF for each language
	for _, lang := range languages {
		logger.Logger().Info("Generating PDF for language", "lang", lang)

		data, err := loader.LoadResumeData(dataDir, lang)
		if err != nil {
			return fmt.Errorf("failed to load resume data for %s: %w", lang, err)
		}

		if err := GeneratePDF(data, outputDir, lang, theme); err != nil {
			return fmt.Errorf("failed to generate PDF for %s: %w", lang, err)
		}
	}

	logger.Logger().Info("PDF generation completed successfully!")
	return nil
}

// GeneratePDF generates a PDF resume for the specified language using the given data and theme.
func GeneratePDF(data *models.ResumeData, outputDir, lang, theme string) error {
	wd, _ := os.Getwd()
	templateDir := filepath.Join(wd, "templates")

	pdfGen, err := generator.NewPDFGenerator(outputDir, templateDir, theme)
	if err != nil {
		return fmt.Errorf("create PDF generator: %w", err)
	}

	return pdfGen.Generate(data, lang)
}
