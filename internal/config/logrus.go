package config

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewLogger(viper *viper.Viper) *logrus.Logger {
	log := logrus.New()

	var writers []io.Writer
	writers = append(writers, os.Stdout)

	file, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		writers = append(writers, file)
	} else {
		logrus.Warnf("Failed to open log file, using stdout only: %v", err)
	}

	log.SetOutput(io.MultiWriter(writers...))

	log.SetLevel(logrus.Level(viper.GetInt("LOG_LEVEL")))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
