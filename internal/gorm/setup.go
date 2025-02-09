package gorm

import (
	"github.com/danbrato999/yuno-gveloz/domain/services"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/models"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/stores"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Order{},
		&models.OrderDish{},
		&models.OrderStatus{},
	)
}

func NewOrderStore(db *gorm.DB) services.OrderStore {
	return stores.NewOrderStore(db)
}

func NewOrderStatusStore(db *gorm.DB) services.OrderStatusStore {
	return stores.NewOrderStatusStore(db)
}
