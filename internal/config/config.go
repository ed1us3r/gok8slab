package config

import (
	"github.com/spf13/viper"
	"fmt"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.gok8slab")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Errorf("❌No config file found. Using defaults.: %v", err)
	}
}

