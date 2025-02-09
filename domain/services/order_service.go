package services

import . "github.com/danbrato999/yuno-gveloz/domain"

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
