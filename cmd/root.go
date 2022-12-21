package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dachshund",
	Short: "dachshund is a web-crawler for everyone",
	Long:  "Use dachshund to crawl your website for quality purposes",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
