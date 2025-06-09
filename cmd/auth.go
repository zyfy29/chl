/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zyfy29/chl/config"
	"github.com/zyfy29/chl/feishu"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := feishu.Api.Auth(config.Conf.Feishu.AppID, config.Conf.Feishu.AppSecret)
		if err != nil {
			cmd.Println("Auth failed:", err)
			return
		}
		if res.Code != 0 {
			cmd.Println("Auth failed:", res.Msg)
			return
		}
		cmd.Println("Auth successful, Tenant Access Token:", res.TenantAccessToken)
		if err := config.SetFeishuTenantAccessToken(res.TenantAccessToken); err != nil {
			cmd.Println("Failed to set Tenant Access Token in config:", err)
			return
		}
		cmd.Println("Tenant Access Token saved to config successfully.")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
