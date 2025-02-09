package models

import "gorm.io/gorm"

type OrderDish struct {
	gorm.Model
	OrderID uint
	Name    string
}
