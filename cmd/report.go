package cmd

import (
	"fmt"
	"log"

	"github.com/connorwade/dachshund/internal"
	"github.com/spf13/cobra"
)

var report = &cobra.Command{
	Use:   "report",
	Short: "Write report from json file",
	Long:  "Write report from json file",
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.WriteCSVReport()
		if err != nil {
			log.Fatalln("Error writing report: ", err)
		}
		fmt.Println("Report has been written")
	},
}

func init() {
	rootCmd.AddCommand(report)
}
