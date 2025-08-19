package config

import (
	"github.com/spf13/viper"
	"log"
)

var Conf Config

type Config struct {
	Base struct {
		TicketID  int    `mapstructure:"ticket_id"`
		GrpcAddr  string `mapstructure:"grpc_addr"`
		Grpc2Addr string `mapstructure:"grpc2_addr"`
	}
	Feishu struct {
		AppID             string `mapstructure:"app_id"`
		AppSecret         string `mapstructure:"app_secret"`
		TenantAccessToken string `mapstructure:"tenant_access_token"`
	} `mapstructure:"feishu"`
	Table struct {
		TableToken string `mapstructure:"table_token"`
		SheetID    string `mapstructure:"sheet_id"`
	}
	Header struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		// chr
		Token   string `mapstructure:"token"`
		Balance string `mapstructure:"balance"`
		Result  string `mapstructure:"result"`
		// gyt
		Cookie      string `mapstructure:"cookie"`
		ItemID      string `mapstructure:"item_id"`
		OrderSN     string `mapstructure:"order_sn"`
		QRCode      string `mapstructure:"qrcode"`
		QRCodeImage string `mapstructure:"qrcode_image"`
	}
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.chl")
	viper.AddConfigPath("..")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}

func setConfig(key, value string) error {
	viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

func SetFeishuTenantAccessToken(token string) error {
	return setConfig("feishu.tenant_access_token", token)
}

func SetSheetID(sheetID string) error {
	return setConfig("table.sheet_id", sheetID)
}
