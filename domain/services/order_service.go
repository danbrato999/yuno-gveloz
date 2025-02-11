package services

import "github.com/danbrato999/yuno-gveloz/domain"

type OrderService interface {
	CreateOrder(request domain.NewOrder) (*domain.Order, error)
	FindByID(id uint) (*domain.OrderWithStatusHistory, error)
	FindMany(filters ...domain.OrderFilterFn) ([]domain.Order, error)
	UpdateStatus(id uint, status domain.OrderStatus) (*domain.Order, error)
	UpdateDishes(id uint, dishes []domain.Dish) (*domain.Order, error)
	Prioritize(id uint, afterID uint) error
}

type orderServiceImpl struct {
	orderStore    OrderStore
	statusStore   OrderStatusStore
	priorityQueue PriorityQueue
}

func NewOrderService(store OrderStore, priorityQueue PriorityQueue, statusStore OrderStatusStore) OrderService {
	return &orderServiceImpl{
		orderStore:    store,
		priorityQueue: priorityQueue,
		statusStore:   statusStore,
	}
}

func (s *orderServiceImpl) CreateOrder(request domain.NewOrder) (*domain.Order, error) {
	order := domain.Order{
		NewOrder: request,
		Status:   domain.OrderStatusPending,
	}

	result, err := s.orderStore.Save(order)
	if err != nil {
		return nil, err
	}

	go s.statusStore.AddCurrentStatus(result)
	go s.priorityQueue.Add(result)

	return result, nil
}

func (s *orderServiceImpl) FindByID(id uint) (*domain.OrderWithStatusHistory, error) {
	order, err := s.findByID(id)

	if err != nil {
		return nil, err
	}

	history, err := s.statusStore.GetHistory(id)
	if err != nil {
		return nil, err
	}

	return &domain.OrderWithStatusHistory{
		Order:         *order,
		StatusHistory: history,
	}, nil
}

func (s *orderServiceImpl) FindMany(filters ...domain.OrderFilterFn) ([]domain.Order, error) {
	orderFilters := &domain.OrderFilters{}

	for _, filter := range filters {
		filter(orderFilters)
	}

	return s.orderStore.GetAll(orderFilters)
}

func (s *orderServiceImpl) UpdateStatus(id uint, status domain.OrderStatus) (*domain.Order, error) {
	existing, err := s.findActiveOrder(id)
	if err != nil {
		return nil, err
	}

	if !existing.IsNewStatusValid(status) {
		return nil, domain.ErrInvalidOrderUpdate
	}

	existing.Status = status

	result, err := s.orderStore.Save(*existing)
	if err != nil {
		return nil, err
	}

	go s.statusStore.AddCurrentStatus(result)

	if status == domain.OrderStatusDone || status == domain.OrderStatusCancelled {
		go s.priorityQueue.Remove(id)
	}

	return result, nil
}

func (s *orderServiceImpl) UpdateDishes(id uint, dishes []domain.Dish) (*domain.Order, error) {
	if len(dishes) == 0 {
		return nil, domain.ErrInvalidOrderUpdate
	}

	existing, err := s.findByID(id)
	if err != nil {
		return nil, err
	}

	if existing.Status != domain.OrderStatusPending && existing.Status != domain.OrderStatusPreparing {
		return nil, domain.ErrInvalidOrderUpdate
	}

	existing.Dishes = dishes

	return s.orderStore.Save(*existing)
}

func (s *orderServiceImpl) Prioritize(id uint, afterID uint) error {
	return s.priorityQueue.ShuffleAfter(id, afterID)
}

func (s *orderServiceImpl) findActiveOrder(id uint) (*domain.Order, error) {
	existing, err := s.findByID(id)

	if err != nil {
		return nil, err
	}

	if existing.Status == domain.OrderStatusDone || existing.Status == domain.OrderStatusCancelled {
		return nil, domain.ErrCompleteOrderUpdate
	}

	return existing, nil
}

func (s *orderServiceImpl) findByID(id uint) (*domain.Order, error) {
	order, err := s.orderStore.FindByID(id)

	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, domain.ErrOrderNotFound
	}

	return order, nil
}
