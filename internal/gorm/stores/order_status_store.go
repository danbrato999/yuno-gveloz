package stores

import (
	"github.com/danbrato999/yuno-gveloz/domain"
	"github.com/danbrato999/yuno-gveloz/domain/services"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/models"
	"gorm.io/gorm"
)

type orderStatusStore struct {
	db *gorm.DB
}

func NewOrderStatusStore(db *gorm.DB) services.OrderStatusStore {
	return &orderStatusStore{
		db: db,
	}
}

// TODO: Check order exists
// TODO: Check current status is not latest
func (o *orderStatusStore) AddCurrentStatus(order *domain.Order) error {
	status := models.OrderStatus{
		OrderID: order.ID,
		Status:  order.Status,
	}

	return o.db.Save(&status).Error
}

func (o *orderStatusStore) GetHistory(id uint) ([]domain.OrderStatusHistory, error) {
	var history []models.OrderStatus

	if err := o.db.Where("order_id = ?", id).Find(&history).Error; err != nil {
		return nil, err
	}

	result := make([]domain.OrderStatusHistory, len(history))

	for i, status := range history {
		result[i] = domain.OrderStatusHistory{
			Status:    status.Status,
			Timestamp: &status.CreatedAt,
		}
	}

	return result, nil
}
