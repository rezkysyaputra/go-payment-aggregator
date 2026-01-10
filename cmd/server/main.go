package main

import (
	"fmt"
	"go-payment-aggregator/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewGin()
	redis := config.NewRedis(viperConfig, log)

	config.Bootstrap(&config.BootstrapConfig{
		Config:   viperConfig,
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Redis:    redis,
	})

	port := viperConfig.GetInt("SERVER_PORT")
	err := app.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
