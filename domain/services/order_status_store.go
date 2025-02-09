package services

import "github.com/danbrato999/yuno-gveloz/domain"

type OrderStatusStore interface {
	AddCurrentStatus(order *domain.Order) error
	GetHistory(id uint) ([]domain.OrderStatusHistory, error)
}
