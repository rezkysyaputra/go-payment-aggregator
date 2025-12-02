package main

import (
	"fmt"
	"go-payment-aggregator/internal/config"
)

func main() {
	// init config
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabse(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewGin()

	// bootstrap
	config.Bootstrap(&config.BootstrapConfig{
		Config:   viperConfig,
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
	})

	// run server
	port := viperConfig.GetInt("server.port")
	err := app.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
