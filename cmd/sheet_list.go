/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zyfy29/chl/config"
	"github.com/zyfy29/chl/feishu"
)

// sheetListCmd represents the list command
var sheetListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sheets, err := feishu.Api.GetSheets(config.Conf.Table.TableToken)
		if err != nil {
			cmd.PrintErrf("Error retrieving sheets: %v\n", err)
			return
		}
		if len(sheets) == 0 {
			cmd.Println("No sheets found.")
			return
		}
		cmd.Printf("index\tID\tTitle\n")
		for _, sheet := range sheets {
			cmd.Printf("%d\t%s\t%s\n", sheet.Index, sheet.SheetId, sheet.Title)
		}
		cmd.Println("Total sheets:", len(sheets))
	},
}

func init() {
	sheetCmd.AddCommand(sheetListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sheetListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sheetListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
