package main

import (
	"github.com/spf13/viper"
	"log"
)

var Conf Config

type Config struct {
	Feishu struct {
		AppID     string `mapstructure:"app_id"`
		AppSecret string `mapstructure:"app_secret"`
	} `mapstructure:"feishu"`
	Table struct {
		TableToken string `mapstructure:"table_token"`
		SheetID    string `mapstructure:"sheet_id"`
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.chl")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	log.Printf("Config loaded: %+v", Conf)
}
