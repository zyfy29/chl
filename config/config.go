package config

import (
	"github.com/spf13/viper"
	"log"
)

var Conf Config

type Config struct {
	Feishu struct {
		AppID             string `mapstructure:"app_id"`
		AppSecret         string `mapstructure:"app_secret"`
		TenantAccessToken string `mapstructure:"tenant_access_token"`
	} `mapstructure:"feishu"`
	Table struct {
		TableToken string `mapstructure:"table_token"`
		SheetID    string `mapstructure:"sheet_id"`
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
