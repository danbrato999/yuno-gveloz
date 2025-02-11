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
	var positions []models.OrderPosition

	return o.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("order_id IN (?, ?)", id, targetID).Find(&positions).Error

		if err != nil {
			return err
		}

		// Any of the orders have no priority do nothing
		if len(positions) < 2 {
			return nil
		}

		currentPos := positions[0].Position
		targetPos := positions[1].Position

		// In case the order gets messed up
		if positions[0].OrderID != id {
			currentPos = positions[1].Position
			targetPos = positions[0].Position
		}

		// We want to add afterwards, not replace
		targetPos++

		if currentPos == targetPos {
			return nil
		}

		if currentPos > targetPos {
			err = tx.Exec("update order_positions set position = position + 1 where position >= ? AND position < ?", targetPos, currentPos).Error

		} else {
			err = tx.Exec("update order_positions set position = position - 1 where position > ? AND position < ?", currentPos, targetPos).Error
			targetPos--
		}

		if err != nil {
			return err
		}

		return tx.Exec("update order_positions set position = ? where order_id = ?", targetPos, id).Error
	})
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
