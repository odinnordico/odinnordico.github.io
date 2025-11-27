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

// WebsiteCmd represents the website command for generating static websites.
var WebsiteCmd = &cobra.Command{
	Use:   "website",
	Short: "Generate static website from YAML resume data",
	Long:  `Generate creates a static website from YAML data files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dataDir := viper.GetString("data-dir")
		outputDir := viper.GetString("output-dir")
		theme := viper.GetString("theme")

		logger.Logger().Info("Starting website generation...")
		logger.Logger().Info("Data directory", "dataDir", dataDir)
		logger.Logger().Info("Output directory", "outputDir", outputDir)
		logger.Logger().Info("Theme", "theme", theme)

		// Validate input directories
		if err := utils.ValidateDirectories(dataDir); err != nil {
			return fmt.Errorf("directory validation failed: %w", err)
		}

		// Ensure output directory exists
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Generate website for all languages
		if err := GenerateMultiLanguageWebsite(dataDir, outputDir, utils.DefaultLang, theme); err != nil {
			return err
		}

		logger.Logger().Info("Website generation completed successfully!")
		return nil
	},
}

func init() {
	WebsiteCmd.Flags().String("theme", "default", "website theme to use")
	viper.BindPFlag("theme", WebsiteCmd.Flags().Lookup("theme"))
}

// GenerateMultiLanguageWebsite generates websites for all available languages.
// The default language (English) is placed in the root output directory,
// while other languages are placed in subdirectories (e.g., /es for Spanish).
func GenerateMultiLanguageWebsite(dataDir, outputDir, defaultLang, theme string) error {
	langDirName := "lang"
	languages := []string{defaultLang}
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

	// Generate website for each language
	for _, lang := range languages {
		logger.Logger().Info("Generating website for", "lang", lang)

		// Determine localized output directory
		localizedOutputDir := outputDir
		if lang != utils.DefaultLang {
			localizedOutputDir = filepath.Join(outputDir, lang)
		}

		// Load resume data
		data, err := loader.LoadResumeData(dataDir, lang)
		if err != nil {
			return fmt.Errorf("load resume data for %s: %w", lang, err)
		}

		// Generate static website (only copy assets for default language)
		copyAssets := lang == utils.DefaultLang
		if err := GenerateWebsite(data, localizedOutputDir, lang, theme, copyAssets); err != nil {
			return fmt.Errorf("generate website for %s: %w", lang, err)
		}
	}

	return nil
}

// GenerateWebsite generates a static website for a single language.
func GenerateWebsite(data *models.ResumeData, outputDir, lang, theme string, copyAssets bool) error {
	wd, _ := os.Getwd()
	templatesDir := filepath.Join(wd, "templates")
	assetsDir := filepath.Join(wd, "assets")

	websiteGen := generator.NewWebsiteGenerator(templatesDir, theme, assetsDir)

	return websiteGen.Generate(data, outputDir, lang, copyAssets)
}
