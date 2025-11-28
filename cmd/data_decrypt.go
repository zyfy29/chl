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

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt tokens based on length and TP balance check",
	Long:  `Decrypt tokens based on length and TP balance check. Tokens longer than 140 chars will always be decrypted, tokens shorter than 140 chars will never be decrypted, tokens exactly 140 chars will be decrypted only if GetTpBalance returns error containing '401004'`,
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
		cmd.Println("Processing tokens for decryption...")

		tokenColIdx := -1
		for i, cell := range res.ValueRange.Data[0] {
			if cell == config.Conf.Header.Token {
				tokenColIdx = i
				break
			}
		}
		if tokenColIdx == -1 {
			cmd.Println("Token column not found in the data.")
			return
		}

		for i, row := range res.ValueRange.Data {
			if i == 0 {
				continue // Skip header row
			}

			token := strings.TrimSpace(row[tokenColIdx])
			if token == "" {
				continue
			}

			originalToken := token
			shouldDecrypt := false

			// Length-based logic
			if len(token) > 140 {
				token = token[:140] // Trim to 140 characters for processing
				shouldDecrypt = true
			} else if len(token) < 140 {
				shouldDecrypt = false
			} else { // exactly 140 characters
				// Call GetTpBalance to check if decryption is needed
				PApi := psdk.NewClient(token, 500, nil, config.Conf.Base.PocketPA)
				_, err := PApi.GetTpBalance(config.Conf.Base.TicketID)
				if err != nil && strings.Contains(err.Error(), "401004") {
					shouldDecrypt = true
				}
			}

			if shouldDecrypt {
				// Decrypt: reverse characters from position 47 to 81 (1-based indexing)
				beforeDecrypt := token[46:81] // positions 47-81 (0-based: 46-80)
				reversed := reverseString(beforeDecrypt)
				decryptedToken := token[:46] + reversed + token[81:]

				// Update the token in the sheet
				if err := feishu.Api.WriteCellData(config.Conf.Table.TableToken, config.Conf.Table.SheetID, feishu.Index2Range(i, tokenColIdx), decryptedToken); err != nil {
					cmd.Printf("Row %d: Error writing decrypted token to cell [%d, %d]: %v\n", i, i, tokenColIdx, err)
				} else {
					cmd.Printf("Row %d: Token decrypted (length: %d)\n", i, len(originalToken))
				}
			} else {
				cmd.Printf("Row %d: Token not decrypted (length: %d)\n", i, len(token))
			}
		}
	},
}

// reverseString reverses a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func init() {
	dataCmd.AddCommand(decryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
