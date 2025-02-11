package gorm

import (
	"fmt"
	"os"
	"time"

	"github.com/danbrato999/yuno-gveloz/domain/services"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/models"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/stores"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DbFolder = "data"

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Order{},
		&models.OrderDish{},
		&models.OrderStatus{},
	)
}

func GetDBConnection(dbName string) (*gorm.DB, error) {
	if err := os.MkdirAll(DbFolder, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error creating db folder: %w", err)
	}

	dbFile := fmt.Sprintf("%s/%s.db?_journal_mode=WAL", DbFolder, dbName)

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	if err = migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database handle: %w", err)
	}

	// Set connection pool limits
	sqlDB.SetMaxOpenConns(1)                  // Max number of open connections
	sqlDB.SetMaxIdleConns(2)                  // Max number of idle connections
	sqlDB.SetConnMaxLifetime(time.Minute * 5) // Reuse connections for 5 minutes

	return db, nil
}

func NewOrderStore(db *gorm.DB) services.OrderStore {
	return stores.NewOrderStore(db)
}

func NewOrderStatusStore(db *gorm.DB) services.OrderStatusStore {
	return stores.NewOrderStatusStore(db)
}
