package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(config *Config) *gorm.DB {
	dsn := config.Dsn
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		panic("Failed to connect database")
	} else {
		println("Database connection established")
	}
	return db
}
