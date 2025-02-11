package stores

import (
	"errors"

	"github.com/danbrato999/yuno-gveloz/domain"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/models"
	"gorm.io/gorm"
)

type OrderPositionStore struct {
	db *gorm.DB
}

func NewOrderPositionStore(db *gorm.DB) *OrderPositionStore {
	return &OrderPositionStore{
		db: db,
	}
}

func (o *OrderPositionStore) Add(order *domain.Order) error {
	err := o.db.
		Where("order_id = ?", order.ID).
		First(new(models.OrderPosition)).
		Error

	if err == nil {
		return domain.ErrIncorrectOrderQueueing
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return o.db.Transaction(func(tx *gorm.DB) error {
		var latest uint
		err := tx.Model(&models.OrderPosition{}).Select("COALESCE(MAX(position), 0)").Scan(&latest).Error
		if err != nil {
			return err
		}

		return tx.Save(&models.OrderPosition{
			OrderID:  order.ID,
			Position: latest + 1,
		}).Error
	})
}

func (o *OrderPositionStore) ShuffleAfter(id, targetID uint) error {
	return nil
}

func (o *OrderPositionStore) Remove(id uint) error {
	return o.db.Transaction(func(tx *gorm.DB) error {
		var current models.OrderPosition

		err := tx.Where("order_id = ?", id).First(&current).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		if err != nil {
			return err
		}

		err = tx.Exec("update order_positions set position = position - 1 where position > ?", current.Position).Error
		if err != nil {
			return err
		}

		return tx.Unscoped().Delete(&current).Error
	})
}
