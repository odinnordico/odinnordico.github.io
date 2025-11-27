// Package cmd provides command-line interface commands for the odinnordico.github.io application.
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/odinnordico/odinnordico.github.io/internal/logger"
	"github.com/odinnordico/odinnordico.github.io/internal/utils"
)

const (
	defaultPort           = "8080"
	defaultHost           = "localhost"
	debounceInterval      = 100 * time.Millisecond
	defaultFilePermission = 0755
)

// ServeCmd represents the serve command for running a local development server.
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the generated website locally",
	Long:  `Serve starts a local HTTP server to preview the generated website for testing and development. Use --watch to enable live reloading.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dataDir := viper.GetString("data-dir")
		outputDir := viper.GetString("output-dir")
		lang := utils.GetLang(viper.GetString("lang"))
		port := viper.GetString("port")
		host := viper.GetString("host")
		watch := viper.GetBool("watch")
		theme := viper.GetString("theme")

		// Create regeneration function
		regenerateWebsite := createRegenerationFunc(dataDir, outputDir, lang, theme)

		// Initial generation if needed
		if err := ensureWebsiteExists(outputDir, watch, regenerateWebsite); err != nil {
			return err
		}

		// Start file watcher if enabled
		if watch {
			if err := startFileWatcher(dataDir, regenerateWebsite); err != nil {
				return err
			}
		}

		// Start HTTP server
		return startHTTPServer(outputDir, host, port, watch)
	},
}

func init() {
	ServeCmd.Flags().String("port", defaultPort, "port to serve on")
	ServeCmd.Flags().String("host", defaultHost, "host to serve on")
	ServeCmd.Flags().Bool("watch", false, "enable live reloading when files change")
	ServeCmd.Flags().String("theme", "default", "website theme to use")

	viper.BindPFlag("port", ServeCmd.Flags().Lookup("port"))
	viper.BindPFlag("host", ServeCmd.Flags().Lookup("host"))
	viper.BindPFlag("watch", ServeCmd.Flags().Lookup("watch"))
	viper.BindPFlag("theme", ServeCmd.Flags().Lookup("theme"))
}

// createRegenerationFunc returns a function that regenerates the website and PDF.
func createRegenerationFunc(dataDir, outputDir, lang, theme string) func() error {
	return func() error {
		logger.Logger().Info("Regenerating website...")

		// Validate directories
		if err := utils.ValidateDirectories(dataDir); err != nil {
			return fmt.Errorf("directory validation failed: %w", err)
		}

		// Ensure output directory exists
		if err := os.MkdirAll(outputDir, defaultFilePermission); err != nil {
			return fmt.Errorf("create output directory: %w", err)
		}

		// Generate website for all languages
		if err := GenerateMultiLanguageWebsite(dataDir, outputDir, lang, theme); err != nil {
			return fmt.Errorf("generate website: %w", err)
		}

		// Generate PDF
		if err := GenerateMultiLanguagePdf(dataDir, outputDir, lang, theme); err != nil {
			return fmt.Errorf("generate PDF: %w", err)
		}

		logger.Logger().Info("Website regenerated successfully!")
		return nil
	}
}

// ensureWebsiteExists checks if the website exists and generates it if needed.
func ensureWebsiteExists(websiteDir string, watch bool, regenerate func() error) error {
	_, err := os.Stat(websiteDir)
	needsGeneration := os.IsNotExist(err) || watch

	if needsGeneration {
		return regenerate()
	}

	if err != nil {
		return fmt.Errorf("website directory does not exist: %s. Please run 'odinnordico.github.io website' first or use --watch to auto-generate", websiteDir)
	}

	return nil
}

// startFileWatcher starts watching for file changes and triggers regeneration.
func startFileWatcher(dataDir string, regenerate func() error) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create file watcher: %w", err)
	}

	// Watch data, templates, and assets directories
	dirs := []string{dataDir, "templates", "assets"}
	for _, dir := range dirs {
		if err := watchDirectory(watcher, dir); err != nil {
			return err
		}
	}

	// Start watching in a goroutine
	go handleFileChanges(watcher, regenerate)

	return nil
}

// watchDirectory recursively adds all subdirectories to the watcher.
func watchDirectory(watcher *fsnotify.Watcher, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
}

// handleFileChanges processes file system events and triggers regeneration.
func handleFileChanges(watcher *fsnotify.Watcher, regenerate func() error) {
	var lastEvent time.Time

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Debounce events to avoid multiple rapid regenerations
			if time.Since(lastEvent) < debounceInterval {
				continue
			}
			lastEvent = time.Now()

			// Only regenerate on write operations
			if event.Has(fsnotify.Write) {
				logger.Logger().Info("File changed", "file", event.Name)
				if err := regenerate(); err != nil {
					logger.Logger().Error("Failed to regenerate website", "error", err)
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Logger().Error("File watcher error", "error", err)
		}
	}
}

// startHTTPServer starts the HTTP server to serve the website.
func startHTTPServer(websiteDir, host, port string, watch bool) error {
	fs := http.FileServer(http.Dir(websiteDir))

	// Handler to serve index.html for root requests
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(websiteDir, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf("%s:%s", host, port)

	fmt.Printf("🚀 Starting local server at http://%s\n", addr)
	fmt.Printf("📁 Serving directory: %s\n", websiteDir)
	if watch {
		fmt.Println("👀 Watching for file changes...")
	}
	fmt.Println("💡 Press Ctrl+C to stop the server")

	return http.ListenAndServe(addr, handler)
}
