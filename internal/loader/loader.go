package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"

	"github.com/odinnordico/odinnordico.github.io/internal/logger"
	"github.com/odinnordico/odinnordico.github.io/internal/models"
)

var (
	extensions     = []string{"yml", "yaml"}
	supportedFiles = map[string]func(*models.ResumeData) any{
		"basic":        func(rd *models.ResumeData) any { return &rd.Basic },
		"professional": func(rd *models.ResumeData) any { return &rd.Professional },
		"certificates": func(rd *models.ResumeData) any { return &rd.Certificates },
		"education":    func(rd *models.ResumeData) any { return &rd.Education },
		"skills":       func(rd *models.ResumeData) any { return &rd.Skills },
		"social":       func(rd *models.ResumeData) any { return &rd.Social },
	}
)

func LoadResumeData(dataDir, lang string) (*models.ResumeData, error) {
	resumeData := &models.ResumeData{}
	for _, ext := range extensions {
		for file, targetFn := range supportedFiles {
			fileName := fmt.Sprintf("%s.%s", file, ext)
			if err := loadYAMLData(dataDir, fileName, lang, resumeData, targetFn); err != nil {
				logger.Logger().Error("failed to load resume data", "language", lang, "file", fileName, "error", err)
				return nil, err
			}
		}
	}
	return resumeData, nil
}

func loadYAMLData(dataDir, file, lang string, target *models.ResumeData, targetFn func(*models.ResumeData) any) error {
	dataPath := filepath.Join(dataDir, file)
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		logger.Logger().Debug("base YAML file does not exist, skipping", "file", dataPath)
		return nil
	}

	logger.Logger().Debug("loading data from YAML file", "file", dataPath)
	targetResume := targetFn(target)
	if err := loadYAMLFile(dataPath, targetResume); err != nil {
		logger.Logger().Error("failed to load data from base YAML file", "file", file, "error", err)
		return err
	}

	if lang == "en" || lang == "" {
		return nil
	}

	dataPath = filepath.Join(dataDir, "lang", lang, file)
	logger.Logger().Debug("loading data from YAML file", "language", lang, "file", dataPath)
	if _, err := os.Stat(dataPath); err == nil {
		// Create a new instance of the target type for language data
		langTarget := &models.ResumeData{}
		targetLangResume := targetFn(langTarget)
		if err := loadYAMLFile(dataPath, targetLangResume); err != nil {
			logger.Logger().Error("failed to load data from lang YAML file", "language", lang, "file", file, "error", err)
			return err
		}

		// Deep merge language data onto the base data (target)
		if err := mergo.Merge(targetResume, targetLangResume, mergo.WithOverride); err != nil {
			logger.Logger().Error("failed to merge language data", "language", lang, "file", file, "error", err)
			return err
		}
	}

	return nil
}

func loadYAMLFile(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal YAML from %s: %w", path, err)
	}

	return nil
}
