package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewRedis(config *viper.Viper, log *logrus.Logger) *redis.Client {
	host := config.GetString("REDIS_HOST")
	port := config.GetInt("REDIS_PORT")
	password := config.GetString("REDIS_PASSWORD")
	db := config.GetInt("REDIS_DB")

	if host == "" {
		host = "localhost"
	}
	if port == 0 {
		port = 6379
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Errorf("Failed to connect to Redis: %v", err)
	}

	return rdb
}
