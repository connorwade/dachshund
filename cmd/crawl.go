package cmd

import (
	"fmt"
	"log"

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
			err := internal.WriteReports(true, true)
			if err != nil {
				log.Fatalln("Error writing report: ", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(crawl)
	crawl.Flags().BoolVarP(&rep, "report", "R", false, "write report after crawl")
}
