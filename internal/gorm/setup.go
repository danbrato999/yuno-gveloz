package gorm

import (
	"fmt"
	"os"

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

	dbFile := fmt.Sprintf("%s/%s.db", DbFolder, dbName)

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	if err = migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func NewOrderStore(db *gorm.DB) services.OrderStore {
	return stores.NewOrderStore(db)
}

func NewOrderStatusStore(db *gorm.DB) services.OrderStatusStore {
	return stores.NewOrderStatusStore(db)
}
