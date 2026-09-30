package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/liyansasongko/bumdes-be/internal/config"
	"github.com/liyansasongko/bumdes-be/internal/entity"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Permission{},
		&entity.Role{},
		&entity.User{},
		&entity.UserLog{},
		&entity.ParkingPrice{},
		&entity.PaymentMethod{},
		&entity.RevenueCategory{},
		&entity.TicketPrice{},
		&entity.TenantType{},
		&entity.AttractionPrice{},
		&entity.Parking{},
		&entity.KiosLocation{},
		&entity.Ticketing{},
		&entity.Attraction{},
		&entity.Tenant{},
		&entity.BillingHistory{},
	)
}
