package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/odinnordico/odinnordico.github.io/cmd"
	"github.com/odinnordico/odinnordico.github.io/internal/logger"
)

var (
	cfgFile string
)

func main() {
	l := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger.SetupLogger(slog.New(l))

	if err := RootCmd.Execute(); err != nil {
		logger.Logger().Error("failed to run command", "error", err)
		os.Exit(1)
	}
}

var RootCmd = &cobra.Command{
	Use:   "odinnordico.github.io",
	Short: "Generate static websites and PDFs from YAML resume data",
	Long: `odinnordico.github.io is a tool that reads YAML files containing resume and portfolio data
and generates a static website and PDF resume that match specified designs.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.odinnordico.github.io.yaml)")
	RootCmd.PersistentFlags().String("data-dir", "data", "directory containing YAML data files")
	RootCmd.PersistentFlags().String("output-dir", "public", "output directory for generated files")

	viper.BindPFlag("data-dir", RootCmd.PersistentFlags().Lookup("data-dir"))
	viper.BindPFlag("output-dir", RootCmd.PersistentFlags().Lookup("output-dir"))

	RootCmd.AddCommand(cmd.PdfCmd)
	RootCmd.AddCommand(cmd.ServeCmd)
	RootCmd.AddCommand(cmd.WebsiteCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".odinnordico.github.io")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logger.Logger().Info("Using config file", "file", viper.ConfigFileUsed())
	}
}
