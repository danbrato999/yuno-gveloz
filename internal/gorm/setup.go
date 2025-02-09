package gorm

import (
	"github.com/danbrato999/yuno-gveloz/domain"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/models"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/stores"
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
