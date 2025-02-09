package stores

import (
	"errors"

	"github.com/danbrato999/yuno-gveloz/internal/adapters/gorm/mappers"
	"github.com/danbrato999/yuno-gveloz/internal/adapters/gorm/models"
	"github.com/danbrato999/yuno-gveloz/internal/domain"
	"gorm.io/gorm"
)

type orderStore struct {
	db *gorm.DB
}

func NewOrderStore(db *gorm.DB) domain.OrderStore {
	return &orderStore{
		db: db,
	}
}

func (o *orderStore) FindByID(id uint) (*domain.Order, error) {
	var order models.Order

	err := o.db.Preload("Dishes").First(&order, id).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return mappers.OrderFromDB(&order), nil
}

func (o *orderStore) GetAll() ([]*domain.Order, error) {
	var orders []models.Order

	if err := o.db.Find(&orders).Error; err != nil {
		return nil, err
	}

	results := make([]*domain.Order, len(orders))

	for i, order := range orders {
		results[i] = mappers.OrderFromDB(&order)
	}

	return results, nil
}

func (o *orderStore) Save(order domain.Order) (*domain.Order, error) {
	dbOrder := mappers.OrderToDB(order)

	err := o.db.Save(&dbOrder).Error

	if err != nil {
		return nil, err
	}

	order.ID = dbOrder.ID
	return &order, nil
}
