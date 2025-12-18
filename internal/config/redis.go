package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewRedis(config *viper.Viper, log *logrus.Logger) *redis.Client {
	host := config.GetString("redis.host")
	port := config.GetInt("redis.port")
	password := config.GetString("redis.password")
	db := config.GetInt("redis.db")

	// Fallback/Default if config missing
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

	// Test connection
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		// Log but don't panic if Redis is optional, but for queue it is critical
		log.Errorf("Failed to connect to Redis: %v", err)
	}

	return rdb
}
