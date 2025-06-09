/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"chl/config"
	"chl/feishu"
	"github.com/spf13/cobra"
	"strconv"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.PrintErrln("Usage: chl sheet use <sheet_index> | <sheet_id>")
			return
		}
		var sheetId string
		if sheetIndex, err := strconv.Atoi(args[0]); err == nil {
			sheets, err := feishu.Api.GetSheets(config.Conf.Table.TableToken)
			if err != nil {
				cmd.PrintErrf("Error retrieving sheets: %v\n", err)
				return
			}
			if sheetIndex < 0 || sheetIndex >= len(sheets) {
				cmd.PrintErrf("Invalid sheet index: %d. There are only %d sheets available.\n", sheetIndex, len(sheets))
				return
			}
			// If the argument is a valid index, use it to get the sheet ID
			sheetId = sheets[sheetIndex].SheetId
		} else {
			// Otherwise, treat it as a sheet ID
			sheetId = args[0]
		}
		if err := config.SetSheetID(sheetId); err != nil {
			cmd.PrintErrf("Error setting sheet ID: %v\n", err)
			return
		}
		cmd.Printf("Sheet ID set to: %s\n", sheetId)
	},
}

func init() {
	sheetCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
