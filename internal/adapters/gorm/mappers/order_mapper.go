package mappers

import (
	"github.com/danbrato999/yuno-gveloz/internal/adapters/gorm/models"
	"github.com/danbrato999/yuno-gveloz/internal/domain"
	"gorm.io/gorm"
)

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
			Source: order.Source,
			Dishes: dishes,
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
		Status: order.Status,
		Source: order.Source,
		Dishes: dishes,
	}

	if order.ID > 0 {
		dbOrder.Model = gorm.Model{
			ID: order.ID,
		}
	}

	return dbOrder
}
