/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/zyfy29/chl/config"
	"github.com/zyfy29/chl/feishu"
	psdk "github.com/zyfy29/pocketgo"
	"strings"

	"github.com/spf13/cobra"
)

// resultCmd represents the result command
var resultCmd = &cobra.Command{
	Use:   "result",
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
		cmd.Println("Checking data integrity...")

		tokenColIdx, resultColIdx := -1, -1
		for i, cell := range res.ValueRange.Data[0] {
			if cell == config.Conf.Header.Token {
				tokenColIdx = i
			} else if cell == config.Conf.Header.Result {
				resultColIdx = i
			}
		}
		if tokenColIdx == -1 || resultColIdx == -1 {
			cmd.Println("Required columns not found in the data.")
			return
		}

		for i, row := range res.ValueRange.Data {
			if i == 0 {
				continue // Skip header row
			}
			result := "Error"
			token := row[tokenColIdx]
			if strings.TrimSpace(token) == "" {
				continue
			}
			PApi := psdk.NewClient(token, 500, nil, config.Conf.Base.PocketPA)
			data, err := PApi.GetHandshakeList(config.Conf.Base.TicketID)
			if err != nil {
				cmd.Printf("Row %d: Error retrieving result for token: %v\n", i, err)
			} else {
				result = "Failed"
				if data.QueueNumber != "" {
					result = data.QueueNumber
				}
			}

			if err := feishu.Api.WriteCellData(config.Conf.Table.TableToken, config.Conf.Table.SheetID, feishu.Index2Range(i, resultColIdx), result); err != nil {
				cmd.Printf("Row %d: Error writing result '%s' to cell [%d, %d]: %v\n", i, result, i, resultColIdx, err)
			} else {
				cmd.Printf("Row %d: Result is '%s'\n", i, result)
			}
		}
	},
}

func init() {
	dataCmd.AddCommand(resultCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resultCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resultCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
