/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zyfy29/chl/config"
	"github.com/zyfy29/chl/feishu"
	"slices"
	"strings"
)

// trimCmd represents the trim command
var trimCmd = &cobra.Command{
	Use:   "trim",
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

		colIdxToTrim := []int{}
		for i, cell := range res.ValueRange.Data[0] {
			if cell == config.Conf.Header.Username || cell == config.Conf.Header.Password {
				cmd.Printf("Column %s (%d) need to be trim\n", cell, i)
				colIdxToTrim = append(colIdxToTrim, i)
			}
		}

		for i := range res.ValueRange.Data {
			if i == 0 {
				continue // Skip header row
			}
			for j, cell := range res.ValueRange.Data[i] {
				if slices.Contains(colIdxToTrim, j) {
					// Trim the cell value
					trimmedCell := strings.TrimSpace(cell)
					if trimmedCell != cell {
						range_ := feishu.Index2Range(i, j)
						err := feishu.Api.WriteCellData(config.Conf.Table.TableToken, config.Conf.Table.SheetID, range_, trimmedCell)
						if err != nil {
							cmd.PrintErrf("Error writing trimmed data to cell [%d, %d]: %v\n", i, j, err)
						} else {
							cmd.Printf("Trimmed cell [%d, %d]: '%s' -> '%s'\n", i, j, cell, trimmedCell)
						}
					}
				}
			}
		}
	},
}

func init() {
	dataCmd.AddCommand(trimCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// trimCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// trimCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
