package services

import "github.com/danbrato999/yuno-gveloz/domain"

type PriorityQueue interface {
	Add(order *domain.Order) error
	ShuffleAfter(id, targetID uint) error
	Remove(id uint) error
}
