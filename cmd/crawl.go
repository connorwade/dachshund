package cmd

import (
	"fmt"

	"github.com/connorwade/dachshund/internal"
	"github.com/spf13/cobra"
)

var rep bool

var crawl = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl the website",
	Long:  "Crawl the website",
	Run: func(cmd *cobra.Command, args []string) {

		internal.Crawl()
		fmt.Println("Crawl has finished")
		if rep {
			internal.WriteContentReport()
			internal.WriteFailureReport()
		}
	},
}

func init() {
	rootCmd.AddCommand(crawl)
	crawl.Flags().BoolVarP(&rep, "report", "r", false, "write report after crawl")
}
