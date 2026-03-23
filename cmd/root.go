package cmd

import (
	"fmt"
	"os"

	"github.com/seanluce/metadata-engine/crawler/internal/config"
	"github.com/seanluce/metadata-engine/crawler/internal/crawler"
	"github.com/spf13/cobra"
)

var cfg config.Config

var rootCmd = &cobra.Command{
	Use:   "crawler",
	Short: "Crawls a file system and stores metadata via the API",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Root == "" {
			return fmt.Errorf("--root is required")
		}
		if cfg.ApiURL == "" {
			cfg.ApiURL = os.Getenv("API_URL")
		}
		if cfg.ApiURL == "" {
			return fmt.Errorf("--api or API_URL env var is required")
		}
		return crawler.Run(cfg)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&cfg.Root, "root", "", "Root path to crawl (required)")
	rootCmd.Flags().StringVar(&cfg.ApiURL, "api", "", "API base URL (or set API_URL env var)")
	rootCmd.Flags().IntVar(&cfg.Workers, "workers", 8, "Number of worker goroutines")
}
