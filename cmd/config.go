package cmd

import (
	"log"
	"strings"

	"github.com/connorwade/dachshund/internal"
	"github.com/spf13/cobra"
)

var config = &cobra.Command{
	Use:   "config",
	Short: "config the website",
	Long:  "config the website",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		if strings.Contains(url, "https://") || strings.Contains(url, "http://") {
			log.Fatalln("Input Error: Only enter the website's domain. No 'https' or 'http'")
		}

		cfg := internal.CrawlerVars{
			StarterURL:     url,
			AllowedDomains: []string{url},
			Colly: struct {
				MaxDepth         int  "yaml:\"maxDepth\""
				Async            bool "yaml:\"async\""
				ParallelRequests int  "yaml:\"parallelRequests\""
			}{
				MaxDepth:         0,
				Async:            true,
				ParallelRequests: 2,
			},
		}

		internal.CreateConfigFile(cfg)
	},
}

func init() {
	rootCmd.AddCommand(config)
}
