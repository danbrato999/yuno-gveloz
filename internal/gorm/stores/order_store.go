package stores

import (
	"errors"

	"github.com/danbrato999/yuno-gveloz/domain"
	"github.com/danbrato999/yuno-gveloz/domain/services"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/mappers"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type orderStore struct {
	db *gorm.DB
}

func NewOrderStore(db *gorm.DB) services.OrderStore {
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

func (o *orderStore) GetAll(filters *domain.OrderFilters) ([]*domain.Order, error) {
	var orders []models.Order

	tx := o.db.Table("orders")

	if filters != nil && len(filters.AnyStatus) > 0 {
		tx.Where("status in (?)", filters.AnyStatus)
	}

	if err := tx.Find(&orders).Error; err != nil {
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

	err := o.db.Transaction(func(tx *gorm.DB) error {
		if err2 := tx.Omit(clause.Associations).Save(&dbOrder).Error; err2 != nil {
			return err2
		}

		if err2 := tx.Model(&dbOrder).Association("Dishes").Replace(dbOrder.Dishes); err2 != nil {
			return err2
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	order.ID = dbOrder.ID
	return &order, nil
}
