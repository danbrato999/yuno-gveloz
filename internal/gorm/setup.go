package gorm

import (
	"fmt"
	"os"
	"time"

	"github.com/danbrato999/yuno-gveloz/domain/services"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/models"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/stores"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Order{},
		&models.OrderDish{},
		&models.OrderPosition{},
		&models.OrderStatus{},
	)
}

func GetDBConnection(dbName string) (*gorm.DB, error) {
	dsn := os.Getenv("POSTGRES_DSN")

	if dsn == "" {
		// Assume we're running locally and postgres is run with docker compose
		dsn = "host=localhost user=postgres password=example dbname=postgres port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
	sqlDB.SetMaxOpenConns(20)                 // Max number of open connections
	sqlDB.SetMaxIdleConns(10)                 // Max number of idle connections
	sqlDB.SetConnMaxLifetime(time.Minute * 5) // Reuse connections for 5 minutes

	return db, nil
}

func NewOrderStore(db *gorm.DB) services.OrderStore {
	return stores.NewOrderStore(db)
}

func NewOrderPriorityStore(db *gorm.DB) services.PriorityQueue {
	return stores.NewOrderPositionStore(db)
}

func NewOrderStatusStore(db *gorm.DB) services.OrderStatusStore {
	return stores.NewOrderStatusStore(db)
}
