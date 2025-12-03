package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewLogger(viper *viper.Viper) *logrus.Logger {
	log := logrus.New()

	file, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	log.SetOutput(file)

	log.SetLevel(logrus.Level(viper.GetInt("log.level")))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
