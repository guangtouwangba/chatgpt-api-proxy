package db

import (
	"chatgpt-api-proxy/config"
	"chatgpt-api-proxy/internal/db/model"
	"chatgpt-api-proxy/pkg/logger"
	"fmt"

	"github.com/pkg/errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetUpDatabase() {
	dbConfig := config.Store.GetDatabaseConfig()
	if dbConfig == nil {
		logger.Panicf("Failed to load database config")
	}
	if dbConfig != nil && dbConfig.Enabled {
		// default to postgres
		dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.DatabaseName, dbConfig.Password)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			logger.Panicf("Failed to connect to database: %v", err)
		}
		DB = db
		err = createUsageDB()
		if err != nil {
			logger.Panicf("Failed to create usage table: %v", err)
		}
		return
	}
	logger.Warnf("Database is not enabled")
}

func GetDB() *gorm.DB {
	return DB
}

func createUsageDB() error {
	if DB != nil {
		err := DB.AutoMigrate(&model.OpenAIUsage{})
		if err != nil {
			return errors.Wrap(err, "failed to create usage table")
		}
	}
	return nil
}
