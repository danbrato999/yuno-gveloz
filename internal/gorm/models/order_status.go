package models

import (
	"github.com/danbrato999/yuno-gveloz/domain"
	"gorm.io/gorm"
)

type OrderStatus struct {
	gorm.Model
	OrderID uint
	Order   Order
	Status  domain.OrderStatus
}
