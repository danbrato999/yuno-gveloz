package gorm

import (
	"github.com/danbrato999/yuno-gveloz/internal/adapters/gorm/models"
	"github.com/danbrato999/yuno-gveloz/internal/adapters/gorm/stores"
	"github.com/danbrato999/yuno-gveloz/internal/domain"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) (err error) {
	if err = db.AutoMigrate(&models.Order{}); err != nil {
		return
	}

	err = db.AutoMigrate(&models.OrderDish{})
	return
}

func NewOrderStore(db *gorm.DB) domain.OrderStore {
	return stores.NewOrderStore(db)
}
