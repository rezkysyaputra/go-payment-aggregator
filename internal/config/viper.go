package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigFile(".env")
	config.SetConfigType("env")
	config.AddConfigPath(".")

	config.AutomaticEnv()

	if err := config.ReadInConfig(); err != nil {
		fmt.Printf("Warning: .env file not found, relying on system environment variables: %v\n", err)
	}

	return config
}
