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

var statusWeights = map[OrderStatus]uint{
	OrderStatusPending:   10,
	OrderStatusActive:    20,
	OrderStatusReady:     30,
	OrderStatusDone:      40,
	OrderStatusCancelled: 40,
}

type OrderSource string

const OrderSourceInPerson OrderSource = "in_person"
const OrderSourceDelivery OrderSource = "delivery"
const OrderSourcePhone OrderSource = "phone"

var ErrOrderNotFound = fmt.Errorf("Order not found")
var ErrInvalidStatusUpdate = fmt.Errorf("Order cannot be updated to provided status")
var ErrCompleteOrderUpdate = fmt.Errorf("Completed order cannot be updated")

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

	return s.store.Save(order)
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

func (s *OrderService) UpdateStatus(id uint, status OrderStatus) (*Order, error) {
	existing, err := s.findActiveOrder(id)
	if err != nil {
		return nil, err
	}

	if !existing.IsNewStatusValid(status) {
		return nil, ErrInvalidStatusUpdate
	}

	existing.Status = status

	return s.store.Save(*existing)
}

// func (s *OrderService) Update(id uint, request OrderUpdate) (*Order, error) {
// 	existing, err := s.FindByID(id)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if (existing.Status == OrderStatusDone || existing.Status == OrderStatusCancelled) && len(existing.Dishes) > 0 {
// 		return nil, fmt.Errorf("Order contents cannot be updated")
// 	}

// 	existing.Dishes = request.Dishes
// 	existing.Status = request.Status

// 	return s.store.Save(*existing)
// }

func (s *OrderService) findActiveOrder(id uint) (*Order, error) {
	existing, err := s.FindByID(id)

	if err != nil {
		return nil, err
	}

	if (existing.Status == OrderStatusDone || existing.Status == OrderStatusCancelled) && len(existing.Dishes) > 0 {
		return nil, ErrCompleteOrderUpdate
	}

	return existing, nil
}
