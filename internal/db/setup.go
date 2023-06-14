package db

import (
	"chatgpt-api-proxy/config"
	"chatgpt-api-proxy/pkg/logger"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetUpDatabase() {
	dbConfig := config.Store.GetDatabaseConfig()
	if dbConfig.Enabled {
		// default to postgres
		dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.DatabaseName, dbConfig.Password)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			logger.Panicf("Failed to connect to database: %v", err)
		}
		DB = db
	}
}

func GetDB() *gorm.DB {
	return DB
}
