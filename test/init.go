package test

import (
	"go-payment-aggregator/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var app *gin.Engine

var db *gorm.DB

var viperConfig *viper.Viper

var log *logrus.Logger

var validate *validator.Validate

var redisClient *redis.Client

func init() {
	viperConfig = config.NewViper()
	log = logrus.New()
	validate = validator.New()
	app = gin.Default()
	db = config.NewDatabase(viperConfig, log)
	redisClient = config.NewRedis(viperConfig, log)

	config.Bootstrap(&config.BootstrapConfig{
		Config:   viperConfig,
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Redis:    redisClient,
	})
}
