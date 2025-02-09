package stores

import (
	"errors"

	"github.com/danbrato999/yuno-gveloz/domain"
	"github.com/danbrato999/yuno-gveloz/domain/services"
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

	return OrderFromDB(&order), nil
}

func (o *orderStore) GetAll(filters *domain.OrderFilters) ([]domain.Order, error) {
	var orders []models.Order

	tx := o.db.Table("orders").Preload("Dishes")

	if filters != nil && len(filters.AnyStatus) > 0 {
		tx.Where("status in (?)", filters.AnyStatus)
	}

	if err := tx.Find(&orders).Error; err != nil {
		return nil, err
	}

	results := make([]domain.Order, len(orders))

	for i, order := range orders {
		results[i] = *OrderFromDB(&order)
	}

	return results, nil
}

func (o *orderStore) Save(order domain.Order) (*domain.Order, error) {
	dbOrder := OrderToDB(order)

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

func OrderFromDB(order *models.Order) *domain.Order {
	if order == nil {
		return nil
	}

	dishes := make([]domain.Dish, len(order.Dishes))

	for i, dish := range order.Dishes {
		dishes[i] = domain.Dish{
			Name: dish.Name,
		}
	}

	return &domain.Order{
		ID:     order.ID,
		Status: order.Status,
		NewOrder: domain.NewOrder{
			Dishes: dishes,
			Source: order.Source,
			Time:   *order.Time,
		},
	}
}

func OrderToDB(order domain.Order) models.Order {
	dishes := make([]models.OrderDish, len(order.Dishes))

	for i, dish := range order.Dishes {
		dishes[i] = models.OrderDish{
			Name: dish.Name,
		}
	}

	dbOrder := models.Order{
		Dishes: dishes,
		Source: order.Source,
		Status: order.Status,
		Time:   &order.Time,
	}

	if order.ID > 0 {
		dbOrder.Model = gorm.Model{
			ID: order.ID,
		}
	}

	return dbOrder
}
