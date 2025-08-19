/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/zyfy29/chl/config"
	"github.com/zyfy29/chl/feishu"
	"github.com/zyfy29/chl/shop"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"
)

// cookieCmd represents the cookie command
var cookieCmd = &cobra.Command{
	Use:   "cookie",
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

		colUsernameIdx, colPasswordIdx, colCookieIdx := -1, -1, -1
		for i, cell := range res.ValueRange.Data[0] {
			if cell == config.Conf.Header.Username {
				colUsernameIdx = i
			} else if cell == config.Conf.Header.Password {
				colPasswordIdx = i
			} else if cell == config.Conf.Header.Cookie {
				colCookieIdx = i
			}
		}
		if colUsernameIdx == -1 || colPasswordIdx == -1 || colCookieIdx == -1 {
			cmd.Println("Required columns (Username, Password, Token) not found in the header row.")
			return
		}

		conn, err := grpc.NewClient(config.Conf.Base.Grpc2Addr, grpc.WithInsecure())
		if err != nil {
			cmd.PrintErrf("Error connecting to gRPC server: %v\n", err)
			return
		}
		defer conn.Close()
		client := shop.NewShopServiceClient(conn)

		for i, row := range res.ValueRange.Data {
			if i == 0 {
				continue // Skip header row
			}

			if row[colCookieIdx] != "" {
				cmd.Printf("Row %d: Cookie already exists. Skipping login.\n", i)
				continue
			}

			username, password := row[colUsernameIdx], row[colPasswordIdx]
			if username == "" || len(password) < 6 {
				cmd.Printf("Row %d: Invalid username or password length.\n", i)
				continue
			}

			resp, err := client.Login(cmd.Context(), &shop.LoginRequest{
				Username: username,
				Password: password,
			})
			if err != nil {
				cmd.Printf("Row %d: Error retrieving token for user '%s': %v\n", i, username, err)
				continue
			}

			if cookie, _ := resp.Headers["Cookie"]; cookie != "" {
				if err := feishu.Api.WriteCellData(config.Conf.Table.TableToken, config.Conf.Table.SheetID, feishu.Index2Range(i, colCookieIdx), cookie); err != nil {
					cmd.Printf("Row %d: Error writing token '%s' to cell [%d, %d]: %v\n", i, cookie, i, colCookieIdx, err)
				} else {
					cmd.Printf("Row %d: Successfully logged in user '%s'\n", i, username)
				}
			} else {
				cmd.Printf("Row %d: No cookie found for user '%s'.\n", i, username)
			}
		}
	},
}

func init() {
	dataCmd.AddCommand(cookieCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cookieCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cookieCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
