package database

import (
	"fmt"
	"log"
	"payment-system/config"
	"payment-system/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	DB = db

	// Auto migrate models
	err = db.AutoMigrate(
		&models.Customer{},
		&models.PaymentIntent{},
		&models.Payment{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
