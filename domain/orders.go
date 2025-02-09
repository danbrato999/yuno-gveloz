package domain

import (
	"time"
)

type NewOrder struct {
	Time   time.Time   `json:"time" binding:"required"`
	Dishes []Dish      `json:"dishes" binding:"required,min=1,dive"`
	Source OrderSource `json:"source" binding:"oneof=in_person delivery phone"`
}

type Order struct {
	ID     uint        `json:"id"`
	Status OrderStatus `json:"status"`
	NewOrder
}

func (o *Order) IsNewStatusValid(status OrderStatus) bool {
	return statusWeights[o.Status] < statusWeights[status]
}

type OrderStore interface {
	Save(order Order) (*Order, error)
	FindByID(id uint) (*Order, error)
	GetAll() ([]*Order, error)
}
