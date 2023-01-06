package cmd

import (
	"fmt"
	"log"

	"github.com/connorwade/dachshund/internal"
	"github.com/connorwade/dachshund/web"
	"github.com/spf13/cobra"
)

var h bool

var report = &cobra.Command{
	Use:   "report",
	Short: "Write report from json file",
	Long:  "Write report from json file",
	Run: func(cmd *cobra.Command, args []string) {

		err := internal.WriteReports(true, true, true, true)
		if err != nil {
			log.Fatalln("Error writing report: ", err)
		}

		if h {
			web.Serve()
		}
		fmt.Println("Report has been written")
	},
}

func init() {
	rootCmd.AddCommand(report)
	report.Flags().BoolVarP(&h, "html", "H", false, "Display report as HTML")
}
