/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/zyfy29/chl/config"
	"github.com/zyfy29/chl/feishu"
	psdk "github.com/zyfy29/pocketgo"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check token and write tp balance to sheet",
	Long:  `Check token and write tp balance to sheet`,
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

		tokenColIdx, balanceColIdx := -1, -1
		for i, cell := range res.ValueRange.Data[0] {
			if cell == config.Conf.Header.Token {
				tokenColIdx = i
			} else if cell == config.Conf.Header.Balance {
				balanceColIdx = i
			}
		}
		if tokenColIdx == -1 || balanceColIdx == -1 {
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
			balance, err := PApi.GetTpBalance(config.Conf.Base.TicketID)
			if err != nil {
				cmd.Printf("Row %d: Error retrieving balance for token: %v\n", i, err)
			} else {
				result = strconv.Itoa(balance)
			}

			if err := feishu.Api.WriteCellData(config.Conf.Table.TableToken, config.Conf.Table.SheetID, feishu.Index2Range(i, balanceColIdx), result); err != nil {
				cmd.Printf("Row %d: Error writing balance '%s' to cell [%d, %d]: %v\n", i, result, i, balanceColIdx, err)
			} else {
				cmd.Printf("Row %d: Balance is '%s'\n", i, result)
			}
		}
	},
}

func init() {
	dataCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
