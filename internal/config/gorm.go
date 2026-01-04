package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	username := viper.GetString("DATABASE_USERNAME")
	password := viper.GetString("DATABASE_PASSWORD")
	host := viper.GetString("DATABASE_HOST")
	port := viper.GetInt("DATABASE_PORT")
	database := viper.GetString("DATABASE_NAME")
	maxIdleConnection := viper.GetInt("DATABASE_POOL_IDLE")
	maxConnection := viper.GetInt("DATABASE_POOL_MAX")
	maxLifetimeConnection := viper.GetInt("DATABASE_POOL_LIFETIME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, username, password, database, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})

	if err != nil {
		log.Fatalf("failed to connect databse: %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect databse: %v", err)
	}

	// set connection pool
	connection.SetConnMaxIdleTime(time.Duration(maxIdleConnection))
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Duration(maxLifetimeConnection))

	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...any) {
	l.Logger.Tracef(message, args...)
}

// migrate -path ./database/migrations -database "postgres://postgres:postgres@localhost:5432/payment_aggregator?sslmode=disable" up
// migrate -path ./database/migrations -database "postgres://postgres:postgres@localhost:5432/payment_aggregator?sslmode=disable" down
