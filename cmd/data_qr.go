/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zyfy29/chl/config"
	"github.com/zyfy29/chl/feishu"
	"github.com/zyfy29/chl/shop"
	"google.golang.org/grpc"
	"strings"
)

// qrCmd represents the qr command
var qrCmd = &cobra.Command{
	Use:   "qr",
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

		cookieColIdx, orderSnColIdx, qrcodeColIdx, qrcodeImageColIdx := -1, -1, -1, -1
		for i, cell := range res.ValueRange.Data[0] {
			if cell == config.Conf.Header.Cookie {
				cookieColIdx = i
			} else if cell == config.Conf.Header.OrderSN {
				orderSnColIdx = i
			} else if cell == config.Conf.Header.QRCode {
				qrcodeColIdx = i
			} else if cell == config.Conf.Header.QRCodeImage {
				qrcodeImageColIdx = i
			}
		}
		if cookieColIdx == -1 || orderSnColIdx == -1 || qrcodeColIdx == -1 || qrcodeImageColIdx == -1 {
			cmd.Println("Required columns not found in the data.")
			return
		}

		conn, err := grpc.NewClient(config.Conf.Base.Grpc2Addr, grpc.WithInsecure())
		if err != nil {
			cmd.PrintErrf("Error connecting to gRPC server: %v\n", err)
			return
		}
		defer conn.Close()

		for i, row := range res.ValueRange.Data {
			if i == 0 {
				continue // Skip header row
			}
			if row[qrcodeColIdx] != "" {
				cmd.Println("Row %d: QRCode already exists")
				continue // Skip rows where Order SN already exists
			}
			token := row[cookieColIdx]
			if strings.TrimSpace(token) == "" {
				continue
			}

			client := shop.NewShopServiceClient(conn)
			res, err := client.GetTicketCode(cmd.Context(), &shop.GetTicketCodeRequest{
				Cookie:  token,
				OrderSn: row[orderSnColIdx],
			})

			if err != nil {
				cmd.Printf("Row %d: Error retrieving result for token: %v\n", i, err)
				continue
			}
			codes := res.TicketCode
			for j, code := range codes {
				if err := feishu.Api.WriteCellData(config.Conf.Table.TableToken, config.Conf.Table.SheetID, feishu.Index2Range(i+j, qrcodeColIdx), code.Qrcode); err != nil {
					cmd.Printf("Row %d: Error writing result '%s' to cell [%d, %d]: %v\n", i+j, code.Qrcode, i+j, orderSnColIdx, err)
				} else {
					cmd.Printf("Row %d: Result is '%s'\n", i, code.Qrcode)
				}
				if err := feishu.Api.WriteCellImage(config.Conf.Table.TableToken, config.Conf.Table.SheetID, feishu.Index2Range(i+j, qrcodeImageColIdx), code.QrcodeImage, row[orderSnColIdx]+"."+code.QrcodeExt); err != nil {
					cmd.Printf("Row %d: Error writing image '%s' to cell [%d, %d]: %v\n", i+j, code.QrcodeImage, i+j, qrcodeImageColIdx, err)
				} else {
					cmd.Printf("Row %d: Image written successfully\n", i+j)
				}
			}
		}
	},
}

func init() {
	dataCmd.AddCommand(qrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// qrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// qrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
