package services

import . "github.com/danbrato999/yuno-gveloz/domain"

type OrderService struct {
	orderStore  OrderStore
	statusStore OrderStatusStore
}

func NewOrderService(store OrderStore, statusStore OrderStatusStore) *OrderService {
	return &OrderService{
		orderStore:  store,
		statusStore: statusStore,
	}
}

func (s *OrderService) CreateOrder(request NewOrder) (*Order, error) {
	order := Order{
		NewOrder: request,
		Status:   OrderStatusPending,
	}

	result, err := s.orderStore.Save(order)
	if err != nil {
		return nil, err
	}

	go s.statusStore.AddCurrentStatus(result)

	return result, nil
}

func (s *OrderService) FindByID(id uint) (*OrderWithStatusHistory, error) {
	order, err := s.findByID(id)

	if err != nil {
		return nil, err
	}

	history, err := s.statusStore.GetHistory(id)
	if err != nil {
		return nil, err
	}

	return &OrderWithStatusHistory{
		Order:         *order,
		StatusHistory: history,
	}, nil
}

func (s *OrderService) FindMany(filters ...OrderFilterFn) ([]Order, error) {
	orderFilters := &OrderFilters{}

	for _, filter := range filters {
		filter(orderFilters)
	}

	return s.orderStore.GetAll(orderFilters)
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

	result, err := s.orderStore.Save(*existing)
	if err != nil {
		return nil, err
	}

	go s.statusStore.AddCurrentStatus(result)

	return result, nil
}

func (s *OrderService) findActiveOrder(id uint) (*Order, error) {
	existing, err := s.findByID(id)

	if err != nil {
		return nil, err
	}

	if existing.Status == OrderStatusDone || existing.Status == OrderStatusCancelled {
		return nil, ErrCompleteOrderUpdate
	}

	return existing, nil
}

func (s *OrderService) findByID(id uint) (*Order, error) {
	order, err := s.orderStore.FindByID(id)

	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, ErrOrderNotFound
	}

	return order, nil
}
