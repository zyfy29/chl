/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zyfy29/chl/config"
	"github.com/zyfy29/chl/feishu"
)

// dataListCmd represents the list command
var dataListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := feishu.Api.ReadRangeData(config.Conf.Table.TableToken, config.Conf.Table.SheetID, "")
		if err != nil {
			cmd.PrintErrf("Error retrieving data: %v\n", err)
			return
		}
		if len(res.ValueRange.Data) == 0 {
			cmd.Println("No data found in the specified range.")
			return
		}
		for i, row := range res.ValueRange.Data {
			cmd.Printf("%d\t", i)
			for _, cell := range row {
				cmd.Printf("%s\t", cell)
			}
			cmd.Println()
		}
		cmd.Printf("Table size: %d * %d\n", len(res.ValueRange.Data), len(res.ValueRange.
			Data[0]))
	},
}

func init() {
	dataCmd.AddCommand(dataListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dataListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dataListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
