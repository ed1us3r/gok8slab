package config

import (
	"github.com/spf13/viper"
	"log"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.gok8slab")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No config file found. Using defaults.")
	}
}

