package config

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewLogger(viper *viper.Viper) *logrus.Logger {
	log := logrus.New()

	// Default to stdout (safe for Docker/Cloud)
	var writers []io.Writer
	writers = append(writers, os.Stdout)

	// Try to open file, but don't crash if it fails (e.g. read-only filesystem)
	file, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		writers = append(writers, file)
	} else {
		// Just print error to stdout, don't panic
		logrus.Warnf("Failed to open log file, using stdout only: %v", err)
	}

	log.SetOutput(io.MultiWriter(writers...))

	log.SetLevel(logrus.Level(viper.GetInt("logging.level")))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
