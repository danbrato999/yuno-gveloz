package services

import "github.com/danbrato999/yuno-gveloz/domain"

type OrderStore interface {
	Save(order domain.Order) (*domain.Order, error)
	FindByID(id uint) (*domain.Order, error)
	GetAll(filters *domain.OrderFilters) ([]*domain.Order, error)
}
