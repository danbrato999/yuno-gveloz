package models

type OrderPosition struct {
	OrderID  uint `gorm:"primaryKey"`
	Position uint
}
