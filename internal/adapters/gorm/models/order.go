package models

import (
	"github.com/danbrato999/yuno-gveloz/internal/domain"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Status domain.OrderStatus
	Source domain.OrderSource
	Dishes []OrderDish
}
