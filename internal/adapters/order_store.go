package adapters

import "github.com/danbrato999/yuno-gveloz/internal/domain"

type orderStore struct {
	counter uint
	data    map[uint]*domain.Order
}

func NewOrderStore() domain.OrderStore {
	data := make(map[uint]*domain.Order)
	return &orderStore{
		counter: 1,
		data:    data,
	}
}

func (o *orderStore) FindByID(id uint) (*domain.Order, error) {
	return o.data[id], nil
}

func (o *orderStore) GetAll() ([]*domain.Order, error) {
	orders := make([]*domain.Order, 0, len(o.data))

	for _, order := range o.data {
		orders = append(orders, order)
	}

	return orders, nil
}

func (o *orderStore) Save(order *domain.Order) (*domain.Order, error) {
	if order.ID == 0 {
		order.ID = o.counter
		o.counter += 1
	}

	o.data[order.ID] = order

	return order, nil
}
