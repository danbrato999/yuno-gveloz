package domain

import (
	"fmt"
	"time"
)

type Dish struct {
	Name string `json:"name" binding:"required"`
}

type OrderStatus string

const OrderStatusPending OrderStatus = "pending"
const OrderStatusActive OrderStatus = "active"
const OrderStatusReady OrderStatus = "ready"
const OrderStatusDone OrderStatus = "done"
const OrderStatusCancelled OrderStatus = "cancelled"

type OrderSource string

const OrderSourceInPerson OrderSource = "in_person"
const OrderSourceDelivery OrderSource = "delivery"
const OrderSourcePhone OrderSource = "phone"

var ErrOrderNotFound = fmt.Errorf("Order not found")

type NewOrder struct {
	Time   *time.Time  `json:"time" binding:"-"`
	Dishes []Dish      `json:"dishes" binding:"required,min=1,dive"`
	Source OrderSource `json:"source" binding:"oneof=in_person delivery phone"`
}

type OrderUpdate struct {
	Dishes []Dish      `json:"dishes" binding:"required,min=1,dive"`
	Status OrderStatus `json:"status" binding:"oneof=active ready done cancelled"`
}

type Order struct {
	ID     uint
	Status OrderStatus
	NewOrder
}

type OrderStore interface {
	Save(order *Order) (*Order, error)
	FindByID(id uint) (*Order, error)
	GetAll() ([]*Order, error)
}

type OrderService struct {
	store OrderStore
}

func NewOrderService(store OrderStore) *OrderService {
	return &OrderService{
		store: store,
	}
}

func (s *OrderService) CreateOrder(request NewOrder) (*Order, error) {
	order := Order{
		NewOrder: request,
		Status:   OrderStatusPending,
	}

	return s.store.Save(&order)
}

func (s *OrderService) FindByID(id uint) (*Order, error) {
	order, err := s.store.FindByID(id)

	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, ErrOrderNotFound
	}

	return order, nil
}

func (s *OrderService) FindMany() ([]*Order, error) {
	return s.store.GetAll()
}

func (s *OrderService) Update(id uint, request OrderUpdate) (*Order, error) {
	existing, err := s.FindByID(id)

	if err != nil {
		return nil, err
	}

	if (existing.Status == OrderStatusDone || existing.Status == OrderStatusCancelled) && len(existing.Dishes) > 0 {
		return nil, fmt.Errorf("Order contents cannot be updated")
	}

	existing.Dishes = request.Dishes
	existing.Status = request.Status

	return s.store.Save(existing)
}

func (s *OrderService) Cancel(id uint) (*Order, error) {
	existing, err := s.FindByID(id)

	if err != nil {
		return nil, err
	}

	existing.Status = OrderStatusCancelled

	return s.store.Save(existing)
}
